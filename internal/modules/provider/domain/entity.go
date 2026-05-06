package domain

import "time"

type AuthConfigData struct {
	APIKey       string   `json:"api_key,omitempty"`
	APISecret    string   `json:"api_secret,omitempty"`
	BearerToken  string   `json:"bearer_token,omitempty"`
	OAuthURL     string   `json:"oauth_url,omitempty"`
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type Provider struct {
	ID                   uint64       `json:"id"`
	Name                 string       `json:"name"`
	ProviderType         ProviderType `json:"provider_type"`
	BaseURL              string       `json:"base_url"`
	AuthType             AuthType     `json:"auth_type"`
	APIKeyCiphertext     string       `json:"-"`
	AuthConfigCiphertext string       `json:"-"`
	Config               string       `json:"config"`
	Enabled              bool         `json:"enabled"`
	HealthStatus         HealthStatus `json:"health_status"`
	LastCheckedAt        *time.Time   `json:"last_checked_at"`
	LastError            string       `json:"last_error"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	Deleted              bool         `json:"deleted"`
}

type ModelConfig struct {
	ID           uint64    `json:"id"`
	ProviderID   uint64    `json:"provider_id"`
	Name         string    `json:"name"`
	Model        string    `json:"model"`
	ModelType    ModelType `json:"model_type"`
	Params       string    `json:"params"`
	Capabilities string    `json:"capabilities"`
	IsDefault    bool      `json:"is_default"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Deleted      bool      `json:"deleted"`
}

type ProviderHealth struct {
	ID            uint64    `json:"id"`
	ProviderID    uint64    `json:"provider_id"`
	ModelConfigID uint64    `json:"model_config_id"`
	Status        string    `json:"status"`
	LatencyMs     int64     `json:"latency_ms"`
	ErrorCode     string    `json:"error_code"`
	ErrorMessage  string    `json:"error_message"`
	CheckedAt     time.Time `json:"checked_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type PageRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}
