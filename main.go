package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
)

// Model represents the UI model
type model struct {
	input       textinput.Model
	output      string
	apiKey      string
	client      *openai.Client
	assistantID string
}

// initialModel sets up the initial UI model with all necessary settings
func initialModel(apiKey, assistantID string) model {
	input := textinput.New()
	input.Placeholder = "Type your question"
	input.Focus()
	input.PromptStyle = lipgloss.NewStyle().Bold(true)
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	client := openai.NewClient(apiKey)

	return model{
		input:       input,
		output:      "",
		apiKey:      apiKey,
		client:      client,
		assistantID: assistantID,
	}
}

// processInput handles the interaction with the OpenAI Assistant API
func (m *model) processInput() {
	ctx := context.Background()
	fmt.Println("Starting processInput")

	// Create a thread for the conversation
	thread, err := m.client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		m.output = "Failed to create thread: " + err.Error()
		fmt.Println(m.output)
		return
	}
	fmt.Println("Thread created")

	// Send the user's question as a message to the thread
	_, err = m.client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: m.input.Value(),
	})
	if err != nil {
		m.output = "Failed to send message: " + err.Error()
		fmt.Println(m.output)
		return
	}
	fmt.Println("Message sent")

	// Start the thread with the assistant
	run, err := m.client.CreateRun(ctx, thread.ID, openai.RunRequest{
		AssistantID: m.assistantID,
	})
	if err != nil {
		m.output = "Failed to start the thread: " + err.Error()
		fmt.Println(m.output)
		return
	}
	fmt.Println("Thread started")

	// Wait for the thread run to complete
	for run.Status != openai.RunStatusCompleted {
		fmt.Println("Waiting for run to complete")
		time.Sleep(5 * time.Second)
		run, err = m.client.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			m.output = "Failed to retrieve run status: " + err.Error()
			fmt.Println(m.output)
			return
		}
	}
	fmt.Println("Run completed")

	// Retrieve all messages from the thread
	msgs, err := m.client.ListMessage(ctx, thread.ID, nil, nil, nil, nil)
	if err != nil {
		m.output = "Failed to retrieve messages: " + err.Error()
		fmt.Println(m.output)
		return
	}

	// Display the last message as the response
	if len(msgs.Messages) > 0 {
		m.output = msgs.Messages[len(msgs.Messages)-1].Content[0].Text.Value
		fmt.Println("Response received:", m.output)
	} else {
		m.output = "No response received"
		fmt.Println(m.output)
	}
}

// Init is called to return the initial command
func (m model) Init() tea.Cmd {
	return nil
}

// Update processes key presses and other messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyEnter:
			go m.processInput() // Process the input asynchronously
			m.input.SetValue("")
			m.input.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the UI to the terminal
func (m model) View() string {
	return fmt.Sprintf("%s\n%s", m.input.View(), m.output)
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("API key not set")
		os.Exit(1)
	}

	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")
	if assistantID == "" {
		fmt.Println("Assistant ID not set")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(apiKey, assistantID))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %s\n", err)
		os.Exit(1)
	}
}
