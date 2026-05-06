package application

import (
	"context"
	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
)

type QueryService struct {
	providerRepo    domain.ProviderRepository
	modelConfigRepo domain.ModelConfigRepository
}

func NewQueryService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
) *QueryService {
	return &QueryService{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
	}
}

func (s *QueryService) GetProvider(ctx context.Context, id uint64) (*ProviderDTO, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	dto := ToProviderDTO(provider)
	dto.APIKeyCiphertext = ""
	dto.AuthConfigCiphertext = ""

	return dto, nil
}

func (s *QueryService) ListProviders(ctx context.Context, query *ProviderPageQuery) (*domain.PageResult[*ProviderDTO], error) {
	pageReq := &domain.PageRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := s.providerRepo.FindPage(ctx, pageReq)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 列表失败")
	}

	dtoList := make([]*ProviderDTO, 0, len(result.Items))
	for _, provider := range result.Items {
		dto := ToProviderDTO(provider)
		dto.APIKeyCiphertext = ""
		dto.AuthConfigCiphertext = ""
		dtoList = append(dtoList, dto)
	}

	return &domain.PageResult[*ProviderDTO]{
		Items:    dtoList,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *QueryService) GetProviderConfigForRuntime(ctx context.Context, providerID uint64) (*ProviderDTO, error) {
	provider, err := s.providerRepo.FindByID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	if !provider.Enabled {
		return nil, errors.ForbiddenErr("Provider 已被禁用")
	}

	return ToProviderDTO(provider), nil
}

func (s *QueryService) GetModelConfig(ctx context.Context, id uint64) (*ModelConfigDTO, error) {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return nil, errors.NotFoundErr("模型配置不存在")
	}

	return ToModelConfigDTO(modelConfig), nil
}

func (s *QueryService) ListModelConfigs(ctx context.Context, query *ModelConfigPageQuery) (*domain.PageResult[*ModelConfigDTO], error) {
	pageReq := &domain.PageRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := s.modelConfigRepo.FindPage(ctx, pageReq)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
	}

	dtoList := make([]*ModelConfigDTO, 0, len(result.Items))
	for _, mc := range result.Items {
		if query.ProviderID == 0 || mc.ProviderID == query.ProviderID {
			dtoList = append(dtoList, ToModelConfigDTO(mc))
		}
	}

	return &domain.PageResult[*ModelConfigDTO]{
		Items:    dtoList,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *QueryService) ListModelConfigsByProvider(ctx context.Context, providerID uint64) ([]*ModelConfigDTO, error) {
	modelConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
	}

	dtoList := make([]*ModelConfigDTO, 0, len(modelConfigs))
	for _, mc := range modelConfigs {
		dtoList = append(dtoList, ToModelConfigDTO(mc))
	}

	return dtoList, nil
}

func (s *QueryService) GetModelConfigForRuntime(ctx context.Context, modelConfigID uint64) (*ModelConfigDTOWithProvider, error) {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, modelConfigID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return nil, errors.NotFoundErr("模型配置不存在")
	}

	provider, err := s.providerRepo.FindByID(ctx, modelConfig.ProviderID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	if !provider.Enabled {
		return nil, errors.ForbiddenErr("Provider 已被禁用")
	}

	if !modelConfig.Enabled {
		return nil, errors.ForbiddenErr("模型配置已被禁用")
	}

	return &ModelConfigDTOWithProvider{
		ModelConfigDTO: ToModelConfigDTO(modelConfig),
		ProviderDTO:    ToProviderDTO(provider),
	}, nil
}

type ModelConfigDTOWithProvider struct {
	*ModelConfigDTO
	*ProviderDTO
}
