package application

import (
	"time"

	"lattice-coding/internal/modules/provider/domain"
)

type ProviderDTO struct {
	ID                   uint64
	Name                 string
	ProviderType         string
	BaseURL              string
	AuthType             string
	APIKeyCiphertext     string
	AuthConfigCiphertext string
	Config               string
	Enabled              bool
	HealthStatus         string
	LastCheckedAt        *time.Time
	LastError            string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ProviderDTOWithHealth struct {
	*ProviderDTO
	Health *ProviderHealthDTO
}

type ModelConfigDTO struct {
	ID           uint64
	ProviderID   uint64
	Name         string
	Model        string
	ModelType    string
	Params       string
	Capabilities string
	IsDefault    bool
	Enabled      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProviderHealthDTO struct {
	ID            uint64
	ProviderID    uint64
	ModelConfigID uint64
	Status        string
	LatencyMs     int64
	ErrorCode     string
	ErrorMessage  string
	CheckedAt     time.Time
	CreatedAt     time.Time
}

func ToProviderDTO(p *domain.Provider) *ProviderDTO {
	if p == nil {
		return nil
	}
	return &ProviderDTO{
		ID:                   p.ID,
		Name:                 p.Name,
		ProviderType:         string(p.ProviderType),
		BaseURL:              p.BaseURL,
		AuthType:             string(p.AuthType),
		APIKeyCiphertext:     p.APIKeyCiphertext,
		AuthConfigCiphertext: p.AuthConfigCiphertext,
		Config:               p.Config,
		Enabled:              p.Enabled,
		HealthStatus:         string(p.HealthStatus),
		LastCheckedAt:        p.LastCheckedAt,
		LastError:            p.LastError,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}

func ToModelConfigDTO(m *domain.ModelConfig) *ModelConfigDTO {
	if m == nil {
		return nil
	}
	return &ModelConfigDTO{
		ID:           m.ID,
		ProviderID:   m.ProviderID,
		Name:         m.Name,
		Model:        m.Model,
		ModelType:    string(m.ModelType),
		Params:       m.Params,
		Capabilities: m.Capabilities,
		IsDefault:    m.IsDefault,
		Enabled:      m.Enabled,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToProviderHealthDTO(h *domain.ProviderHealth) *ProviderHealthDTO {
	if h == nil {
		return nil
	}
	return &ProviderHealthDTO{
		ID:            h.ID,
		ProviderID:    h.ProviderID,
		ModelConfigID: h.ModelConfigID,
		Status:        h.Status,
		LatencyMs:     h.LatencyMs,
		ErrorCode:     h.ErrorCode,
		ErrorMessage:  h.ErrorMessage,
		CheckedAt:     h.CheckedAt,
		CreatedAt:     h.CreatedAt,
	}
}
