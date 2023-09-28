package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/j-dunham/openai-cli/services/openai"
	"github.com/joho/godotenv"
	"github.com/muesli/reflow/wordwrap"
)

func main() {
	godotenv.Load()
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	spinner       spinner.Model
	loading       bool
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	responseStyle lipgloss.Style
	err           error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ta := textarea.New()
	ta.Placeholder = "What is your Prompt?"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 200

	ta.SetWidth(50)
	ta.SetHeight(2)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(100, 10)
	vp.SetContent(`Welcome to the OpenAI CLI!
Type a prompt and press ENTER.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		spinner:       s,
		loading:       false,
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		responseStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		err:           nil,
	}
}

func chatCompletion(prompt string) tea.Msg {
	response := openai.GetCompletion(prompt)
	wrapped := wordwrap.String(response, 50)
    return completionMsg(wrapped)
}

type completionMsg string

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
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
		case tea.KeyEnter:
			wrappedPrompt := wordwrap.String(m.textarea.Value(), 50)
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+wrappedPrompt)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			prompt := m.textarea.Value()
			m.textarea.Reset()
			m.viewport.GotoBottom()
			m.loading = true
			return m, func() tea.Msg { return chatCompletion(prompt) }
		}
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

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if m.loading {
		loading := fmt.Sprintf("\n\n   %s Loading...\n\n", m.spinner.View())
		return loading
	}
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
