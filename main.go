package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/j-dunham/openai-cli/config"
	"github.com/j-dunham/openai-cli/services/openai"
	"github.com/j-dunham/openai-cli/services/storage"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

func initialModel(cfg *config.Config) model {
	return model{
		cfg:           cfg,
		spinner:       newSpinner(),
		loading:       false,
		viewport:      newViewport(),
		messages:      []string{},
		textarea:      newTextarea(),
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		responseStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
		openAiService: openai.NewService(cfg),
		table:         newTable(),
		showTable:     false,
		help:          newHelp(),
	}
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel(cfg))
	storage.CreateTable()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	cfg           *config.Config
	spinner       spinner.Model
	loading       bool
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	responseStyle lipgloss.Style
	err           error
	openAiService openai.Service
	table         table.Model
	showTable     bool
	help          string
}

func newTable() table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Prompt", Width: 50},
		{Title: "Response", Width: 50},
	}

	prompts, _ := storage.ReadPrompts()
	rows := make([]table.Row, 0)
	for _, p := range prompts {
		rows = append(rows, table.Row{p.ID, p.Prompt, p.Response})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	tStyle := table.DefaultStyles()
	tStyle.Header = tStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tStyle.Selected = tStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(tStyle)
	return t
}

func newTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "What is your Prompt?"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 200

	ta.SetWidth(50)
	ta.SetHeight(2)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	return ta
}

func newViewport() viewport.Model {
	vp := viewport.New(100, 10)
	vp.SetContent(`Welcome to the OpenAI CLI!
Type a prompt and press ENTER.`)
	return vp
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func newHelp() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	return helpStyle.Render("CTRL+T History | CTRL+C Exit")
}

func savePrompt(prompt string, response string) {
	storage.InsertPrompt(prompt, response)
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

type completionMsg string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmd   tea.Cmd
	)

	blueText := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyCtrlT:
			prompts, _ := storage.ReadPrompts()
			rows := make([]table.Row, 0)
			for _, p := range prompts {
				rows = append(rows, table.Row{p.ID, p.Prompt, p.Response})
			}
			m.table.SetRows(rows)
			m.showTable = !m.showTable
		case tea.KeyEnter:

			if m.showTable {
				wrappedPrompt := wrap.String(m.table.SelectedRow()[1], 50)
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+blueText.Render(wrappedPrompt)+"\n")
				wrappedResponse := wrap.String(m.table.SelectedRow()[2], 50)
				m.messages = append(m.messages, m.responseStyle.Render("OpenAI: ")+blueText.Render(wrappedResponse)+"\n")
			} else {
				wrappedPrompt := wrap.String(m.textarea.Value(), 50)
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+blueText.Render(wrappedPrompt)+"\n")
			}

			prompt := m.textarea.Value()
			if !m.showTable {
				cmd = func() tea.Msg {
					_, msg := m.openAiService.GetCompletion(prompt)
					return msg
				}
				m.loading = true
			} else {
				m.showTable = false
				cmd = nil
			}
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			m.loading = true
			return m, func() tea.Msg {
				response, err := m.openAiService.GetCompletion(prompt)
				if err != nil {
					// not sure if this is how to best handle this error
					// double-check the docs
					return errMsg(err)
				}
				wrapped := wordwrap.String(response, 50)
				return completionMsg(wrapped)
			}
		}
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case completionMsg:
		m.messages = append(m.messages, m.responseStyle.Render("OpenAI: ")+blueText.Render(string(msg))+"\n")
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		m.loading = false
	case errMsg:
		m.err = msg
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(tiCmd, vpCmd, cmd)
}

func (m model) View() string {
	if m.loading {
		loading := fmt.Sprintf("\n\n   %s Loading...\n\n", m.spinner.View())
		return loading
	}
	if m.showTable {
		return fmt.Sprintf(
			"%s\n%s\n%s",
			"Prompt History Selector",
			m.table.View(),
			m.help,
		) + "\n\n"
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
		m.help,
	) + "\n\n"
}
