package application

import (
	"context"
	"fmt"
	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
	"lattice-coding/internal/runtime/llm"
)

type SyncService struct {
	providerRepo    domain.ProviderRepository
	modelConfigRepo domain.ModelConfigRepository
	modelLister     llm.ModelLister
}

func NewSyncService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	modelLister llm.ModelLister,
) *SyncService {
	return &SyncService{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
		modelLister:     modelLister,
	}
}

func (s *SyncService) SyncModels(ctx context.Context, providerID uint64, force bool) (*SyncModelsResponse, error) {
	provider, err := s.providerRepo.FindByID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	if !provider.Enabled {
		return nil, errors.ForbiddenErr("Provider 已被禁用，无法同步模型")
	}

	apiKey := ""
	if provider.APIKeyCiphertext != "" {
		apiKey = provider.APIKeyCiphertext
	}

	modelsResp, err := s.modelLister.ListModels(ctx, string(provider.ProviderType), provider.BaseURL, apiKey)
	if err != nil {
		if err == llm.ErrUnsupportedProviderType {
			return nil, errors.InvalidArg("该 Provider 类型不支持模型列表同步")
		}
		return nil, errors.InternalWithErr(err, "拉取模型列表失败")
	}

	existingConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询现有模型配置失败")
	}

	existingMap := make(map[string]*domain.ModelConfig)
	for _, cfg := range existingConfigs {
		existingMap[cfg.Model] = cfg
	}

	response := &SyncModelsResponse{
		ProviderID: providerID,
		Total:      len(modelsResp.Data),
		Created:    0,
		Skipped:    0,
		Failed:     0,
	}

	for _, model := range modelsResp.Data {
		if _, exists := existingMap[model.ID]; exists {
			response.Skipped++
			continue
		}

		modelConfig := &domain.ModelConfig{
			ProviderID:   providerID,
			Name:         model.ID,
			Model:        model.ID,
			ModelType:    domain.ModelTypeChat,
			Params:       "{}",
			Capabilities: "",
			IsDefault:    false,
			Enabled:      true,
		}

		if err := s.modelConfigRepo.Create(ctx, modelConfig); err != nil {
			response.Failed++
			response.Message += "Failed to create model: " + model.ID + "; "
			continue
		}
		response.Created++
	}

	if response.Created > 0 || response.Failed > 0 {
		response.Message = fmt.Sprintf("同步完成: 共 %d 个模型，成功创建 %d 个，跳过 %d 个，失败 %d 个",
			response.Total, response.Created, response.Skipped, response.Failed)
	} else {
		response.Message = fmt.Sprintf("同步完成: 共 %d 个模型，无需更新", response.Total)
	}

	return response, nil
}

type SyncModelsResponse struct {
	ProviderID uint64 `json:"provider_id"`
	Total      int    `json:"total"`
	Created    int    `json:"created"`
	Skipped    int    `json:"skipped"`
	Failed     int    `json:"failed"`
	Message    string `json:"message"`
}
