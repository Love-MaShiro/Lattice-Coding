package application

import (
	"lattice-coding/internal/modules/provider/domain"
)

func AssembleProviderDTOFromEntity(p *domain.Provider) *ProviderDTO {
	return ToProviderDTO(p)
}

func AssembleModelConfigDTOFromEntity(m *domain.ModelConfig) *ModelConfigDTO {
	return ToModelConfigDTO(m)
}

func AssembleProviderHealthDTOFromEntity(h *domain.ProviderHealth) *ProviderHealthDTO {
	return ToProviderHealthDTO(h)
}
