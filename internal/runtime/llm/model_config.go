package llm

import "context"

type ModelConfigResolver interface {
	ResolveModelConfig(ctx context.Context, modelConfigID string) (*ResolvedModelConfig, error)
	ResolveAgentDefaultModel(ctx context.Context, agentID string) (*ResolvedModelConfig, error)
}

type ProviderDefaultModelResolver interface {
	ResolveProviderDefaultModel(ctx context.Context, providerID string) (*ResolvedModelConfig, error)
}

type ResolvedModelConfig struct {
	ProviderID    string
	ProviderType  string
	ModelConfigID string
	ModelName     string
	BaseURL       string
	APIKey        string
	Temperature   float64
	TopP          float64
	MaxTokens     int
	Extra         map[string]interface{}
}
