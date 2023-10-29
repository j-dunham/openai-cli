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
)

func initialModel(cfg *config.Config) model {
	return model{
		cfg:           cfg,
		spinner:       newSpinner(),
		loading:       false,
		viewport:      newViewport(),
		messages:      []openai.Message{},
		textarea:      newTextarea(),
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
	messages      []openai.Message
	textarea      textarea.Model
	err           error
	openAiService openai.Service
	table         table.Model
	showTable     bool
	help          string
}

func newTable() table.Model {
	storage.CreateTable()
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
	return helpStyle.Render("CTRL+T History Table  | CTRL+W Wipe History | CTRL+C Exit")
}

func savePrompt(prompt string, response string) {
	storage.InsertPrompt(prompt, response)
}

func RenderMessages(messages []openai.Message) string {
	colors := map[string]lipgloss.Style{
		"user":      lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		"assistant": lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		"system":    lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		"prompt":    lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
	}

	formatedMsgs := make([]string, 0)
	for _, msg := range messages {
		s := colors[msg.Role].Render(strings.ToUpper(msg.Role)) + ": " + colors["prompt"].Render(msg.Content)
		if msg.Role == "assistant" {
			s += "\n"
		}
		formatedMsgs = append(formatedMsgs, wordwrap.String(s, 50))
	}
	return strings.Join(formatedMsgs, "\n")
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
		case tea.KeyCtrlW:
			m.messages = []openai.Message{}
			m.viewport.SetContent(RenderMessages(m.messages))
			return m, cmd
		case tea.KeyEnter:
			if m.showTable {
				m.messages = append(m.messages, openai.Message{Role: "user", Content: m.table.SelectedRow()[1]})
				m.messages = append(m.messages, openai.Message{Role: "assistant", Content: m.table.SelectedRow()[2]})
			} else {
				prompt := m.textarea.Value()
				m.messages = append(m.messages, openai.Message{Role: "user", Content: prompt})
			}
			m.viewport.SetContent(RenderMessages(m.messages))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			if m.showTable {
				m.showTable = false
				cmd = nil
				return m, cmd
			}

			m.loading = true
			return m, func() tea.Msg {
				response, err := m.openAiService.GetCompletion(m.messages)
				if err != nil {
					// not sure if this is how to best handle this error
					// double-check the docs
					return errMsg(err)
				}
				savePrompt(m.messages[len(m.messages)-1].Content, response)
				return completionMsg(response)
			}
		}
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case completionMsg:
		m.messages = append(m.messages, openai.Message{Role: "assistant", Content: string(msg)})
		m.textarea.Reset()
		m.viewport.SetContent(RenderMessages(m.messages))
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
