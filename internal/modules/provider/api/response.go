package api

import (
	"time"

	"lattice-coding/internal/modules/provider/application"
	"lattice-coding/internal/modules/provider/domain"
)

type ProviderResponse struct {
	ID            uint64     `json:"id"`
	Name          string     `json:"name"`
	ProviderType  string     `json:"provider_type"`
	BaseURL       string     `json:"base_url"`
	AuthType      string     `json:"auth_type"`
	APIKeySet     bool       `json:"api_key_set"`
	Config        string     `json:"config"`
	Enabled       bool       `json:"enabled"`
	HealthStatus  string     `json:"health_status"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	LastError     string     `json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ProviderPageResponse struct {
	Items    []ProviderResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type ModelConfigResponse struct {
	ID           uint64    `json:"id"`
	ProviderID   uint64    `json:"provider_id"`
	ProviderName string    `json:"provider_name,omitempty"`
	Name         string    `json:"name"`
	Model        string    `json:"model"`
	ModelType    string    `json:"model_type"`
	Params       string    `json:"params,omitempty"`
	Capabilities string    `json:"capabilities,omitempty"`
	IsDefault    bool      `json:"is_default"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ModelConfigPageResponse struct {
	Items    []ModelConfigResponse `json:"items"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}

type ProviderHealthResponse struct {
	ProviderID    uint64    `json:"provider_id"`
	ModelConfigID uint64    `json:"model_config_id,omitempty"`
	Status        string    `json:"status"`
	LatencyMs     int64     `json:"latency_ms"`
	ErrorCode     string    `json:"error_code,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	CheckedAt     time.Time `json:"checked_at"`
}

type ProviderTestResponse struct {
	Success   bool   `json:"success"`
	LatencyMs int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

type ModelTestResponse struct {
	Success   bool   `json:"success"`
	LatencyMs int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

func ToProviderResponse(dto *application.ProviderDTO) ProviderResponse {
	if dto == nil {
		return ProviderResponse{}
	}
	return ProviderResponse{
		ID:            dto.ID,
		Name:          dto.Name,
		ProviderType:  dto.ProviderType,
		BaseURL:       dto.BaseURL,
		AuthType:      dto.AuthType,
		APIKeySet:     dto.APIKeyCiphertext != "",
		Config:        dto.Config,
		Enabled:       dto.Enabled,
		HealthStatus:  dto.HealthStatus,
		LastCheckedAt: dto.LastCheckedAt,
		LastError:     dto.LastError,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func ToModelConfigResponse(dto *application.ModelConfigDTO) ModelConfigResponse {
	if dto == nil {
		return ModelConfigResponse{}
	}
	return ModelConfigResponse{
		ID:           dto.ID,
		ProviderID:   dto.ProviderID,
		Name:         dto.Name,
		Model:        dto.Model,
		ModelType:    dto.ModelType,
		Params:       dto.Params,
		Capabilities: dto.Capabilities,
		IsDefault:    dto.IsDefault,
		Enabled:      dto.Enabled,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func ToProviderHealthResponse(dto *application.ProviderHealthDTO) ProviderHealthResponse {
	if dto == nil {
		return ProviderHealthResponse{}
	}
	return ProviderHealthResponse{
		ProviderID:    dto.ProviderID,
		ModelConfigID: dto.ModelConfigID,
		Status:        dto.Status,
		LatencyMs:     dto.LatencyMs,
		ErrorCode:     dto.ErrorCode,
		ErrorMessage:  dto.ErrorMessage,
		CheckedAt:     dto.CheckedAt,
	}
}

func ToProviderPageResponse(result *domain.PageResult[*application.ProviderDTO]) *ProviderPageResponse {
	if result == nil {
		return nil
	}
	items := make([]ProviderResponse, 0, len(result.Items))
	for _, dto := range result.Items {
		items = append(items, ToProviderResponse(dto))
	}
	return &ProviderPageResponse{
		Items:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}
}

func ToModelConfigPageResponse(result *domain.PageResult[*application.ModelConfigDTO]) *ModelConfigPageResponse {
	if result == nil {
		return nil
	}
	items := make([]ModelConfigResponse, 0, len(result.Items))
	for _, dto := range result.Items {
		items = append(items, ToModelConfigResponse(dto))
	}
	return &ModelConfigPageResponse{
		Items:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}
}
