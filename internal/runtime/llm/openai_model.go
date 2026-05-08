package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

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
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		modelName:  modelName,
		httpClient: &http.Client{},
	}, nil
}

var _ model.ChatModel = (*OpenAIChatModel)(nil)

func (m *OpenAIChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	reqBody := m.buildRequest(messages, false, opts...)
	respBytes, err := m.doJSON(ctx, reqBody)
	if err != nil {
		return nil, err
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
	reqBody := m.buildRequest(messages, true, opts...)
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	m.setHeaders(req)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBytes))
	}

	sr, sw := schema.Pipe[*schema.Message](16)
	go func() {
		defer resp.Body.Close()
		defer sw.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, ":") {
				continue
			}
			if !strings.HasPrefix(line, "data:") {
				continue
			}

			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "[DONE]" {
				return
			}

			var chunk openAIStreamResponse
			if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
				if sw.Send(nil, err) {
					return
				}
				continue
			}
			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta.Content
			if delta == "" {
				continue
			}
			if sw.Send(&schema.Message{Role: schema.Assistant, Content: delta}, nil) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			sw.Send(nil, err)
		}
	}()

	return sr, nil
}

func (m *OpenAIChatModel) BindTools(tools []*schema.ToolInfo) error {
	return nil
}

func (m *OpenAIChatModel) buildRequest(messages []*schema.Message, stream bool, opts ...model.Option) *openAIRequest {
	options := model.GetCommonOptions(nil, opts...)
	modelName := m.modelName
	if options.Model != nil && *options.Model != "" {
		modelName = *options.Model
	}

	reqBody := &openAIRequest{
		Model:    modelName,
		Stream:   stream,
		Messages: make([]openAIMessage, len(messages)),
	}
	if options.Temperature != nil {
		reqBody.Temperature = options.Temperature
	}
	if options.TopP != nil {
		reqBody.TopP = options.TopP
	}
	if options.MaxTokens != nil {
		reqBody.MaxTokens = options.MaxTokens
	}
	if len(options.Stop) > 0 {
		reqBody.Stop = options.Stop
	}

	for i, msg := range messages {
		reqBody.Messages[i] = openAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	return reqBody
}

func (m *OpenAIChatModel) doJSON(ctx context.Context, reqBody *openAIRequest) ([]byte, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	m.setHeaders(req)

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
	return respBytes, nil
}

func (m *OpenAIChatModel) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Stream      bool            `json:"stream,omitempty"`
	Temperature *float32        `json:"temperature,omitempty"`
	TopP        *float32        `json:"top_p,omitempty"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
	Stop        []string        `json:"stop,omitempty"`
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

type openAIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}
