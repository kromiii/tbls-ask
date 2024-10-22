package client

import (
	"context"
	"testing"
)

type MockChatClient struct{
	ChatCompletionRequest(ctx context.Context, prompt string) (string, error)
}

func TestNew(t *testing.T) {
	tests := []struct{
		name string
		model string
		querymode bool
		want *Client[ChatClientType]
		wantErr bool
	}{
		{
			name: "openai chat client",
			model: "gpt-3.5-turbo",
			querymode: false,
			want: &Client[ChatClientType]{
				chatClient: "gpt-3.5-turbo",
				querymode: false,
				promptTmpl: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New[ChatClientType](context.Background(), tt.model, tt.querymode)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.chatClient != tt.want.chatClient {
				t.Errorf("New() = %v, want %v", got.chatClient, tt.want.chatClient)
			}
			if got.querymode != tt.want.querymode {
				t.Errorf("New() = %v, want %v", got.querymode, tt.want.querymode)
			}
		})
	}
}
