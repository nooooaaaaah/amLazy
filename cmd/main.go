package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nooooaaaaah/amLazy/config"
	"github.com/nooooaaaaah/amLazy/openai"
	"github.com/nooooaaaaah/amLazy/tui"
)

func main() {
	config.LoadEnv()
	config.InitLogger()
	defer config.GetLogger().Close()

	logger := config.GetLogger()
	logger.LogInfo("Starting amLazy application")
	apiKey := config.GetEnv("OPENAI_API_KEY")
	assistantID := config.GetEnv("OPENAI_ASSISTANT_ID")
	if apiKey == "" || assistantID == "" {
		logger.LogError("API key or Assistant ID not set")
		os.Exit(1)
	}

	client := openai.NewClient(apiKey, assistantID)
	model := tui.InitialModel(client)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		logger.LogErrorf("Error running program: %s\n", err)
		os.Exit(1)
	}
}
