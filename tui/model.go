package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nooooaaaaah/amLazy/config"
	"github.com/nooooaaaaah/amLazy/openai"
)

type Model struct {
	Input  textinput.Model
	Output string
	Client *openai.Client
	Sent   bool
}

func InitialModel(client *openai.Client) *Model {
	input := textinput.New()
	input.Placeholder = "Type your question"
	input.Focus()
	input.PromptStyle = lipgloss.NewStyle().Bold(true)
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	return &Model{
		Input:  input,
		Output: "",
		Client: client,
		Sent:   false,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.Sent { // Only process the input if it hasn't been sent
				m.Sent = true  // Set the flag to true after sending
				m.Input.Blur() // Optionally blur the input after sending
				cmd = m.makeAPICall()
			}
			return m, cmd
		}
	case string:
		// Handle the response message from makeAPICall
		m.Output = msg
		return m, nil
	}

	m.Input, cmd = m.Input.Update(msg) // Continue handling other inputs
	return m, cmd
}

// makeAPICall handles the interaction with OpenAI asynchronously and updates the model
func (m *Model) makeAPICall() tea.Cmd {
	return func() tea.Msg {
		logger := config.GetLogger()
		logger.LogInfo("Sending input to OpenAI: " + m.Input.Value())
		response, err := m.Client.ProcessInput(m.Input.Value())
		if err != nil {
			logger.LogError("Error processing input: " + err.Error())
			return "Error: " + err.Error()
		}
		return response
	}
}

func (m *Model) View() string {
	if m.Sent && m.Output == "" {
		return fmt.Sprintf("%s\nSending...", m.Input.View())
	}
	return fmt.Sprintf("%s\n%s", m.Input.View(), m.Output)
}
