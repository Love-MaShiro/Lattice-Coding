package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

type ModelLister interface {
	ListModels(ctx context.Context, providerType string, baseURL string, apiKey string) (*ModelsResponse, error)
}

func NewModelLister() ModelLister {
	return &defaultModelLister{}
}

type defaultModelLister struct{}

func (l *defaultModelLister) ListModels(ctx context.Context, providerType string, baseURL string, apiKey string) (*ModelsResponse, error) {
	switch providerType {
	case "openai", "openai_compatible":
		return listOpenAICompatibleModels(ctx, baseURL, apiKey)
	case "ollama":
		return listOllamaModels(ctx, baseURL)
	default:
		return nil, ErrUnsupportedProviderType
	}
}

func listOpenAICompatibleModels(ctx context.Context, baseURL string, apiKey string) (*ModelsResponse, error) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	url := baseURL + "/models"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New("failed to list models: status " + string(rune(resp.StatusCode)) + ", body: " + string(bodyBytes))
	}

	var result struct {
		Object string `json:"object"`
		Data   []struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			Created int64  `json:"created"`
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]ModelInfo, len(result.Data))
	for i, m := range result.Data {
		models[i] = ModelInfo{
			ID:      m.ID,
			Object:  m.Object,
			Created: m.Created,
			OwnedBy: m.OwnedBy,
		}
	}

	return &ModelsResponse{
		Object: result.Object,
		Data:   models,
	}, nil
}

type ollamaTagsResponse struct {
	Models []struct {
		Name       string `json:"name"`
		Model      string `json:"model"`
		Size       int64  `json:"size"`
		ModifiedAt string `json:"modified_at"`
	} `json:"models"`
}

func listOllamaModels(ctx context.Context, baseURL string) (*ModelsResponse, error) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	url := baseURL + "/api/tags"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to list ollama models: status " + string(rune(resp.StatusCode)))
	}

	var result ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]ModelInfo, len(result.Models))
	for i, m := range result.Models {
		models[i] = ModelInfo{
			ID:      m.Name,
			Object:  "model",
			Created: 0,
			OwnedBy: "local",
		}
	}

	return &ModelsResponse{
		Object: "list",
		Data:   models,
	}, nil
}
