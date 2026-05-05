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

type OpenAIChatModel struct {
	baseURL    string
	apiKey     string
	modelName  string
	httpClient *http.Client
}

func NewOpenAIChatModel(baseURL, apiKey, modelName string) (*OpenAIChatModel, error) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIChatModel{
		baseURL:    baseURL,
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{},
	}, nil
}

var _ model.ChatModel = (*OpenAIChatModel)(nil)

func (m *OpenAIChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	reqBody := &openAIRequest{
		Model:    m.modelName,
		Messages: make([]openAIMessage, len(messages)),
	}
	for i, msg := range messages {
		reqBody.Messages[i] = openAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := m.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}

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

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBytes, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, errors.New("no response from model")
	}

	return &schema.Message{
		Role:    schema.Assistant,
		Content: openAIResp.Choices[0].Message.Content,
	}, nil
}

func (m *OpenAIChatModel) Stream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("streaming not implemented")
}

func (m *OpenAIChatModel) BindTools(tools []*schema.ToolInfo) error {
	return nil
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}
