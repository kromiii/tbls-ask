// llm.go
// ほしい機能
// 1. クライアントを作成する
// OpenAI または Gemini のモデルを指定してクエリクライアントを作成する
// 2. プロンプトを作成する
// tbls schema の情報を使ってLLMに渡すプロンプトを作成する
// 3. LLM に問い合わせる
// 先ほど初期化したクライアントを使ってLLMに問い合わせる
package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/repin"
	"github.com/k1LoW/tbls/schema"
	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls-ask/internal/openai"
	"github.com/k1LoW/tbls-ask/internal/gemini"
)

const (
	DefaultModelChat = "gpt-4o"
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type LLMClient struct {
	client *openai.OpenAIClient
	model  string
	querymode bool
}

func NewLLMClient(model string, querymode bool) (*LLMClient, error) {
	client, err := openai.NewOpenAIClient(model)
	if err != nil {
		return nil, err
	}
	return &LLMClient{
		client: client,
		model:  model,
		querymode: querymode,
	}, nil
}

func (a *LLMClient) GeneratePrompt(db *schema.Database) (string, error) {
	tmpl, err := template.New("prompt").Parse(templates.PromptTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, db)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *LLMClient) Ask(ctx context.Context, prompt string) (string, error) {
	resp, err := a.client.ChatCompletionRequest(prompt)
	if err != nil {
		return "", err
	}
	return resp, nil
}
