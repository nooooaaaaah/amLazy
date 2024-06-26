package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nooooaaaaah/amLazy/config"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	apiClient   *openai.Client
	assistantID string
}

type MsgRsponse struct {
	Msg string `json:"msg"`
}

// NewClient creates a new Client instance with the specified API key and assistant ID.
func NewClient(apiKey, assistantID string) *Client {
	return &Client{
		apiClient:   openai.NewClient(apiKey),
		assistantID: assistantID,
	}
}

func (c *Client) ProcessInput(input, userEnvInstructions string) (string, error) {
	logger := config.GetLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	logger.LogInfof("processing input: %s", input)

	threadID, err := c.createThread(ctx)
	if err != nil {
		return "", err
	}

	if err := c.sendMessage(ctx, threadID, input); err != nil {
		return "", err
	}

	if err := c.startThread(ctx, threadID, userEnvInstructions); err != nil {
		return "", err
	}

	return c.retrieveResponse(ctx, threadID)
}

func (c *Client) createThread(ctx context.Context) (string, error) {
	logger := config.GetLogger()
	thread, err := c.apiClient.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		logger.LogErrorf("failed to create thread: %v", err)
		return "", fmt.Errorf("failed to create thread: %w", err)
	}
	logger.LogInfo("created a thread")
	return thread.ID, nil
}

func (c *Client) sendMessage(ctx context.Context, threadID string, input string) error {
	logger := config.GetLogger()
	_, err := c.apiClient.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: input,
	})
	if err != nil {
		logger.LogErrorf("failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}
	logger.LogInfo("created a message")
	return nil
}

func (c *Client) startThread(ctx context.Context, threadID, userEnvInstructions string) error {
	logger := config.GetLogger()
	run, err := c.apiClient.CreateRun(ctx, threadID, openai.RunRequest{AssistantID: c.assistantID, Instructions: userEnvInstructions})
	if err != nil {
		logger.LogErrorf("failed to start the thread: %v", err)
		return fmt.Errorf("failed to start the thread: %w", err)
	}

	logger.LogInfo("started the thread")
	for run.Status != openai.RunStatusCompleted {
		select {
		case <-ctx.Done():
			return fmt.Errorf("process timed out or cancelled")
		case <-time.After(5 * time.Second):
			run, err = c.apiClient.RetrieveRun(ctx, threadID, run.ID)
			if err != nil {
				logger.LogErrorf("failed to retrieve run status: %v", err)
				return fmt.Errorf("failed to retrieve run status: %w", err)
			}
			if run.Status == openai.RunStatusCompleted {
				break
			}
		}
	}
	return nil
}

func (c *Client) retrieveResponse(ctx context.Context, threadID string) (string, error) {
	logger := config.GetLogger()
	msgs, err := c.apiClient.ListMessage(ctx, threadID, nil, nil, nil, nil)
	if err != nil {
		logger.LogErrorf("failed to retrieve messages: %v", err)
		return "", fmt.Errorf("failed to retrieve messages: %w", err)
	}

	logger.LogInfo("retrieved messages")
	if len(msgs.Messages) > 0 {
		logger.LogInfof("%v", msgs.Messages)
		response := msgs.Messages[0].Content[0].Text.Value
		var msg MsgRsponse
		err := json.Unmarshal([]byte(response), &msg)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal response message: %w", err)
		}
		logger.LogInfof("First response: %s", response)
		logger.LogInfo(msg.Msg)
		return msg.Msg, nil
	}
	return "No response received", nil
}
