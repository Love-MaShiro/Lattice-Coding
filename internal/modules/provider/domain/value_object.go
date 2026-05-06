package domain

type ProviderType string

const (
	ProviderTypeOpenAI           ProviderType = "openai"
	ProviderTypeOpenAICompatible ProviderType = "openai_compatible"
	ProviderTypeOllama           ProviderType = "ollama"
	ProviderTypeClaude           ProviderType = "claude"
)

type AuthType string

const (
	AuthTypeAPIKey       AuthType = "api_key"
	AuthTypeAPIKeySecret AuthType = "api_key_secret"
	AuthTypeBearerToken  AuthType = "bearer_token"
	AuthTypeOAuth        AuthType = "oauth"
	AuthTypeNone         AuthType = "none"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
	HealthStatusDegraded  HealthStatus = "degraded"
)

type ModelType string

const (
	ModelTypeChat       ModelType = "chat"
	ModelTypeCompletion ModelType = "completion"
	ModelTypeEmbedding  ModelType = "embedding"
	ModelTypeImage      ModelType = "image"
	ModelTypeAudio      ModelType = "audio"
)
