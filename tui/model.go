package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Input  textinput.Model
	Output string
}

func InitialModel() Model {
	input := textinput.New()
	input.Placeholder = "Type your question"
	input.Focus()
	input.PromptStyle = lipgloss.NewStyle().Bold(true)
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	return Model{
		Input:  input,
		Output: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.Input.View(), m.Output)
}
