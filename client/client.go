package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/repin"
	"github.com/k1LoW/tbls-ask/internal/gemini"
	"github.com/k1LoW/tbls-ask/internal/openai"
	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/schema"
)

const (
	DefaultModelChat = "gpt-4o"
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type ChatClientType interface {
	ChatCompletionRequest(ctx context.Context, prompt string) (string, error)
}

type Client[T ChatClientType] struct {
	chatClient  T
	querymode   bool
	promptTmpl  string
}

func New[T ChatClientType](ctx context.Context, model string, querymode bool) (*Client[T], error) {
	if model == "" {
		model = DefaultModelChat
	}
	var promptTmpl string
	if querymode {
		promptTmpl = templates.DefaultQueryPromptTmpl
	} else {
		promptTmpl = templates.DefaultPromptTmpl
	}

	var c T
	if strings.HasPrefix(model, "gpt") {
			c = any(openai.New(model)).(T)
	} else if strings.HasPrefix(model, "gemini") {
			c = any(gemini.New(ctx, model)).(T)
	} else {
			return nil, fmt.Errorf("unsupported model: %s", model)
	}

	return &Client[T]{
		chatClient: c,
		querymode:  querymode,
		promptTmpl: promptTmpl,
	}, nil
}

func (c *Client[T]) Ask(ctx context.Context, q string, s *schema.Schema) (string, error) {
	p, err := c.GeneratePrompt(s, q)
	if err != nil {
		return "", err
	}
	resp, err := c.chatClient.ChatCompletionRequest(ctx, p)
	if err != nil {
		return "", err
	}
	if c.querymode {
		resp, err = extractQuery(resp)
		if err != nil {
			return "", err
		}
	}
	return resp, nil
}

func (c *Client[T]) GeneratePrompt(s *schema.Schema, q string) (string, error) {
	tpl, err := template.New("").Parse(c.promptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(s),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(s),
		"Viewpoints":      templates.GenerateViewPoints(s),
		"Question":        q,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func extractQuery(resp string) (string, error) {
	if !strings.Contains(resp, quoteStart) {
		return "", fmt.Errorf("failed to pick query from answer: %s", resp)
	}
	query := new(bytes.Buffer)
	if _, err := repin.Pick(strings.NewReader(resp), quoteStart, quoteEnd, true, query); err != nil {
		return "", fmt.Errorf("failed to pick query from answer: %w\nanswer: %s", err, resp)
	}
	return query.String(), nil
}
