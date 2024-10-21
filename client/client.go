package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/repin"
	"github.com/k1LoW/tbls-ask/internal/openai"
	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/schema"
)

const (
	DefaultModelChat = "gpt-4o"
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type Client struct {
	chatClient  *openai.OpenAI
	querymode bool
	promptTmpl string
}

func New(model string, querymode bool) (*Client, error) {
	var promptTmpl string
	if querymode {
		promptTmpl = templates.DefaultQueryPromptTmpl
	} else {
		promptTmpl = templates.DefaultPromptTmpl
	}

	c := openai.New(model)

	return &Client{
		chatClient:  c,
		querymode: querymode,
		promptTmpl: promptTmpl,
	}, nil
}

func (c *Client) Ask(ctx context.Context, q string, s *schema.Schema) (string, error) {
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

func (c *Client) GeneratePrompt(s *schema.Schema, q string) (string, error) {
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
