package chat

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

type VertexAIClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewVertexAIClient(model string) (*VertexAIClient, error) {
	ctx := context.Background()

	jsonCredentials := []byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"))
    
    creds, err := google.CredentialsFromJSON(ctx, jsonCredentials, "https://www.googleapis.com/auth/cloud-platform")
    if err != nil {
        return nil, err
    }

	projectID := creds.ProjectID
	if projectID == "" {
		return nil, fmt.Errorf("failed to get project ID from credentials")
	}

	client, err := genai.NewClient(
		ctx,
		projectID,
		"",
		option.WithCredentials(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &VertexAIClient{
		client: client,
		model:  client.GenerativeModel(model),
	}, nil
}

func (c *VertexAIClient) Ask(ctx context.Context, messages []Message) (string, error) {
	chat := c.model.StartChat()

	// Convert messages to Gemini format
	history := make([]*genai.Content, len(messages))
	for i, msg := range messages {
		role := msg.Role
		if role == "system" {
			role = "user"
		} else if role == "assistant" {
			role = "model"
		}

		history[i] = &genai.Content{
			Parts: []genai.Part{
				genai.Text(msg.Content),
			},
			Role: role,
		}
	}

	chat.History = history

	// Send the last message
	lastMsg := messages[len(messages)-1]
	resp, err := chat.SendMessage(ctx, genai.Text(lastMsg.Content))
	if err != nil {
		return "", fmt.Errorf("gemini api error: %w", err)
	}

	// Extract response
	var answer string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			for _, part := range candidate.Content.Parts {
				if part != nil {
					answer = fmt.Sprintf("%s", part)
				}
			}
		}
	}
	return answer, nil
}
