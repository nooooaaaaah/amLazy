package main

import (
	"amLazy/config"
	"amLazy/openai"
	"amLazy/tui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/yourusername/myassistant/config"
	// "github.com/yourusername/myassistant/openai"
	// "github.com/yourusername/myassistant/tui"
)

func main() {
	config.LoadEnv()
	apiKey := config.GetEnv("OPENAI_API_KEY")
	assistantID := config.GetEnv("OPENAI_ASSISTANT_ID")
	if apiKey == "" || assistantID == "" {
		fmt.Println("API key or Assistant ID not set")
		os.Exit(1)
	}

	openai.NewClient(apiKey, assistantID)
	// client := openai.NewClient(apiKey, assistantID)

	model := tui.InitialModel()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %s\n", err)
		os.Exit(1)
	}
}
