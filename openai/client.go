package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/nooooaaaaah/amLazy/config"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	apiClient   *openai.Client
	assistantID string
}

// NewClient creates a new Client instance with the specified API key and assistant ID.
func NewClient(apiKey, assistantID string) *Client {
	return &Client{
		apiClient:   openai.NewClient(apiKey),
		assistantID: assistantID,
	}
}

// ProcessInput handles sending an input to the OpenAI API and retrieving the response.
func (c *Client) ProcessInput(input string) (string, error) {
	logger := config.GetLogger()
	ctx := context.Background()

	// Set a timeout for the API operation
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	logger.LogInfof("processing input: %s", input)

	// Create a new thread
	thread, err := c.apiClient.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		logger.LogErrorf("failed to create thread: %v", err)
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	logger.LogInfo("created a thread")

	// Send the user's input as a message in the thread
	_, err = c.apiClient.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: input,
	})
	if err != nil {
		logger.LogErrorf("failed to send message: %v", err)
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	logger.LogInfo("created a message")

	// Start the thread with the specified assistant
	run, err := c.apiClient.CreateRun(ctx, thread.ID, openai.RunRequest{AssistantID: c.assistantID})
	if err != nil {
		logger.LogErrorf("failed to start the thread: %v", err)
		return "", fmt.Errorf("failed to start the thread: %w", err)
	}

	logger.LogInfo("started the thread")

	// Wait for the thread to complete or timeout
	for run.Status != openai.RunStatusCompleted {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("process timed out or cancelled")
		case <-time.After(5 * time.Second):
			run, err = c.apiClient.RetrieveRun(ctx, thread.ID, run.ID)
			if err != nil {
				logger.LogErrorf("failed to retrieve run status: %v", err)
				return "", fmt.Errorf("failed to retrieve run status: %w", err)
			}
			if run.Status == openai.RunStatusCompleted {
				break
			}
		}
	}

	// Retrieve messages from the thread
	msgs, err := c.apiClient.ListMessage(ctx, thread.ID, nil, nil, nil, nil)
	if err != nil {
		logger.LogErrorf("failed to retrieve messages: %v", err)
		return "", fmt.Errorf("failed to retrieve messages: %w", err)
	}

	logger.LogInfo("retrieved messages")
	if len(msgs.Messages) > 0 {
		// Returning the first message response
		response := msgs.Messages[0].Content[0].Text.Value
		logger.LogInfof("First response: %s", response)
		return response, nil
	}
	return "No response received", nil
}
