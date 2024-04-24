package openai

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	apiClient   *openai.Client
	assistantID string
}

func NewClient(apiKey, assistantID string) *Client {
	return &Client{
		apiClient:   openai.NewClient(apiKey),
		assistantID: assistantID,
	}
}

func (c *Client) ProcessInput(input string) (string, error) {
	ctx := context.Background()

	thread, err := c.apiClient.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = c.apiClient.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: input,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	run, err := c.apiClient.CreateRun(ctx, thread.ID, openai.RunRequest{AssistantID: c.assistantID})
	if err != nil {
		return "", fmt.Errorf("failed to start the thread: %w", err)
	}

	for run.Status != openai.RunStatusCompleted {
		time.Sleep(5 * time.Second)
		run, err = c.apiClient.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve run status: %w", err)
		}
	}

	msgs, err := c.apiClient.ListMessage(ctx, thread.ID, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve messages: %w", err)
	}

	if len(msgs.Messages) > 0 {
		return msgs.Messages[len(msgs.Messages)-1].Content[0].Text.Value, nil
	}
	return "No response received", nil
}
