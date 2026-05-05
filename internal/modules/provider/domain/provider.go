package domain

import "time"

// ProviderType 模型提供商类型
type ProviderType string

const (
	ProviderTypeOpenAI           ProviderType = "openai"
	ProviderTypeOpenAICompatible ProviderType = "openai_compatible"
	ProviderTypeOllama           ProviderType = "ollama"
	ProviderTypeClaude           ProviderType = "claude" // 预留
)

// Provider 模型提供商
type Provider struct {
	ID               uint64       `json:"id"`
	Name             string       `json:"name"`
	ProviderType     ProviderType `json:"provider_type"`
	BaseURL          string       `json:"base_url,omitempty"`
	APIKeyCiphertext string       `json:"-"` // API Key 密文，不返回给前端
	IsEnabled        bool         `json:"is_enabled"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	ID          uint64    `json:"id"`
	ProviderID  uint64    `json:"provider_id"`
	Name        string    `json:"name"`
	ModelName   string    `json:"model_name"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
	ExtraConfig string    `json:"extra_config,omitempty"` // JSON 扩展配置
	IsEnabled   bool      `json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProviderRepository provider 仓库接口
type ProviderRepository interface {
	Create(provider *Provider) error
	Update(provider *Provider) error
	GetByID(id uint64) (*Provider, error)
	List() ([]*Provider, error)
	Delete(id uint64) error
}

// ModelConfigRepository model_config 仓库接口
type ModelConfigRepository interface {
	Create(modelConfig *ModelConfig) error
	GetByID(id uint64) (*ModelConfig, error)
	List() ([]*ModelConfig, error)
	ListByProviderID(providerID uint64) ([]*ModelConfig, error)
	Delete(id uint64) error
}

// CryptoService 加密服务接口（预留）
type CryptoService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
