package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
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
	Info   string
}

var controlKeysStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#acacac")).
	PaddingTop(1).
	PaddingBottom(1)

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
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// config.GetLogger().LogInfof("Key pressed: %v", msg.String()) // Confirm key detection in logs

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.Sent {
				m.Sent = true
				m.Input.Blur()
				cmd = m.makeAPICall()
			}
			return m, cmd
		case tea.KeyCtrlY:
			// Check that the output exists and the API call has been completed
			config.GetLogger().LogInfof("Ctrl+Y pressed: %v", m.Sent) // Ensure this log appears
			if m.Sent && m.Output != "" {
				cmd = copyToClipboard(m.Output)
				return m, cmd
			}
		}

	case string:
		if msg == "copied" {
			m.Info = "Output copied to clipboard!"
			return m, tea.Quit
		}
		m.Output = msg
		return m, nil

	case error:
		m.Info = fmt.Sprintf("Failed to copy: %v", msg)
		return m, nil
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		config.GetLogger().LogInfo("Attempting to copy to clipboard") // New log entry
		if err := clipboard.WriteAll(text); err != nil {
			config.GetLogger().LogErrorf("Copy to clipboard failed: %v", err)
			return fmt.Errorf("copy to clipboard failed: %v", err)
		}
		config.GetLogger().LogInfo("Successfully copied to clipboard") // Confirm successful copy
		return "copied"                                                // Return a success message
	}
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
	controlKeys := controlKeysStyle.Render("[Ctrl+C] Quit [Enter] Send [Ctrl+Y] Copy to Clipboard")
	if m.Sent && m.Output == "" {
		return fmt.Sprintf("%s\nSending...\n%s", m.Input.View(), controlKeys)
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.Input.View(), m.Output, m.Info, controlKeys)
}
