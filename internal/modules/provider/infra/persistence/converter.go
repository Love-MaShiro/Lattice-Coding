package persistence

import (
	"lattice-coding/internal/modules/provider/domain"
)

func ConvertProviderToPO(src *domain.Provider, dst *ProviderPO) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.ProviderType = string(src.ProviderType)
	dst.BaseURL = src.BaseURL
	dst.AuthType = string(src.AuthType)
	dst.APIKeyCiphertext = src.APIKeyCiphertext
	dst.AuthConfigCiphertext = src.AuthConfigCiphertext
	dst.Config = src.Config
	dst.Enabled = src.Enabled
	dst.HealthStatus = string(src.HealthStatus)
	dst.LastCheckedAt = src.LastCheckedAt
	dst.LastError = src.LastError
}

func ConvertPOToProvider(src *ProviderPO) *domain.Provider {
	return &domain.Provider{
		ID:                  src.ID,
		Name:                src.Name,
		ProviderType:        domain.ProviderType(src.ProviderType),
		BaseURL:             src.BaseURL,
		AuthType:            domain.AuthType(src.AuthType),
		APIKeyCiphertext:    src.APIKeyCiphertext,
		AuthConfigCiphertext: src.AuthConfigCiphertext,
		Config:              src.Config,
		Enabled:             src.Enabled,
		HealthStatus:        domain.HealthStatus(src.HealthStatus),
		LastCheckedAt:       src.LastCheckedAt,
		LastError:           src.LastError,
		CreatedAt:           src.CreatedAt,
		UpdatedAt:           src.UpdatedAt,
	}
}

func ConvertModelConfigToPO(src *domain.ModelConfig, dst *ModelConfigPO) {
	dst.ID = src.ID
	dst.ProviderID = src.ProviderID
	dst.Name = src.Name
	dst.Model = src.Model
	dst.ModelType = string(src.ModelType)
	dst.Params = src.Params
	dst.Capabilities = src.Capabilities
	dst.IsDefault = src.IsDefault
	dst.Enabled = src.Enabled
}

func ConvertPOToModelConfig(src *ModelConfigPO) *domain.ModelConfig {
	return &domain.ModelConfig{
		ID:           src.ID,
		ProviderID:   src.ProviderID,
		Name:         src.Name,
		Model:        src.Model,
		ModelType:    domain.ModelType(src.ModelType),
		Params:       src.Params,
		Capabilities: src.Capabilities,
		IsDefault:    src.IsDefault,
		Enabled:      src.Enabled,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func ConvertProviderHealthToPO(src *domain.ProviderHealth, dst *ProviderHealthPO) {
	dst.ID = src.ID
	dst.ProviderID = src.ProviderID
	dst.ModelConfigID = src.ModelConfigID
	dst.Status = src.Status
	dst.LatencyMs = src.LatencyMs
	dst.ErrorCode = src.ErrorCode
	dst.ErrorMessage = src.ErrorMessage
	dst.CheckedAt = src.CheckedAt
}

func ConvertPOToProviderHealth(src *ProviderHealthPO) *domain.ProviderHealth {
	return &domain.ProviderHealth{
		ID:            src.ID,
		ProviderID:    src.ProviderID,
		ModelConfigID: src.ModelConfigID,
		Status:        src.Status,
		LatencyMs:     src.LatencyMs,
		ErrorCode:     src.ErrorCode,
		ErrorMessage:  src.ErrorMessage,
		CheckedAt:     src.CheckedAt,
		CreatedAt:     src.CreatedAt,
	}
}
