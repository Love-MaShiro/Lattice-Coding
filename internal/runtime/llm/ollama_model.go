package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type OllamaChatModel struct {
	baseURL    string
	modelName  string
	httpClient *http.Client
}

func NewOllamaChatModel(baseURL, modelName string) (*OllamaChatModel, error) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaChatModel{
		baseURL:    baseURL,
		modelName:  modelName,
		httpClient: &http.Client{},
	}, nil
}

var _ model.ChatModel = (*OllamaChatModel)(nil)

func (m *OllamaChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	reqBody := &ollamaRequest{
		Model:    m.modelName,
		Messages: make([]ollamaMessage, len(messages)),
		Stream:   false,
	}
	for i, msg := range messages {
		reqBody.Messages[i] = ollamaMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := m.baseURL + "/api/chat"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBytes))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return nil, err
	}

	if ollamaResp.Message.Content == "" {
		return nil, errors.New("no response from model")
	}

	return &schema.Message{
		Role:    schema.Assistant,
		Content: ollamaResp.Message.Content,
	}, nil
}

func (m *OllamaChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("streaming not implemented")
}

func (m *OllamaChatModel) BindTools(tools []*schema.ToolInfo) error {
	return nil
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}
