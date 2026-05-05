package application

import (
	"context"
	"lattice-coding/internal/modules/provider/domain"
)

// CreateProviderRequest 创建 provider 请求
type CreateProviderRequest struct {
	Name         string               `json:"name"`
	ProviderType domain.ProviderType  `json:"provider_type"`
	BaseURL      string               `json:"base_url"`
	APIKey       string               `json:"api_key"` // 前端传明文，我们加密
}

// UpdateProviderRequest 更新 provider 请求
type UpdateProviderRequest struct {
	Name         string               `json:"name"`
	ProviderType domain.ProviderType  `json:"provider_type"`
	BaseURL      string               `json:"base_url"`
	APIKey       string               `json:"api_key"` // 可选，如果不传则不更新
}

// CreateModelConfigRequest 创建 model_config 请求
type CreateModelConfigRequest struct {
	ProviderID  uint64   `json:"provider_id"`
	Name        string   `json:"name"`
	ModelName   string   `json:"model_name"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
	ExtraConfig string   `json:"extra_config,omitempty"`
}

// ProviderService provider 用例服务
type ProviderService struct {
	providerRepo    domain.ProviderRepository
	modelConfigRepo domain.ModelConfigRepository
	cryptoService   domain.CryptoService
}

func NewProviderService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	cryptoService domain.CryptoService,
) *ProviderService {
	return &ProviderService{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
		cryptoService:   cryptoService,
	}
}

func (s *ProviderService) CreateProvider(ctx context.Context, req *CreateProviderRequest) (*domain.Provider, error) {
	apiKeyCiphertext, err := s.cryptoService.Encrypt(req.APIKey)
	if err != nil {
		return nil, err
	}

	provider := &domain.Provider{
		Name:             req.Name,
		ProviderType:    req.ProviderType,
		BaseURL:          req.BaseURL,
		APIKeyCiphertext: apiKeyCiphertext,
		IsEnabled:        true,
	}

	if err := s.providerRepo.Create(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *ProviderService) UpdateProvider(ctx context.Context, id uint64, req *UpdateProviderRequest) (*domain.Provider, error) {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	provider.Name = req.Name
	provider.ProviderType = req.ProviderType
	provider.BaseURL = req.BaseURL

	if req.APIKey != "" {
		apiKeyCiphertext, err := s.cryptoService.Encrypt(req.APIKey)
		if err != nil {
			return nil, err
		}
		provider.APIKeyCiphertext = apiKeyCiphertext
	}

	if err := s.providerRepo.Update(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *ProviderService) GetProvider(ctx context.Context, id uint64) (*domain.Provider, error) {
	return s.providerRepo.GetByID(id)
}

func (s *ProviderService) ListProviders(ctx context.Context) ([]*domain.Provider, error) {
	return s.providerRepo.List()
}

func (s *ProviderService) DeleteProvider(ctx context.Context, id uint64) error {
	return s.providerRepo.Delete(id)
}

func (s *ProviderService) EnableProvider(ctx context.Context, id uint64) error {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return err
	}
	provider.IsEnabled = true
	return s.providerRepo.Update(provider)
}

func (s *ProviderService) DisableProvider(ctx context.Context, id uint64) error {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return err
	}
	provider.IsEnabled = false
	return s.providerRepo.Update(provider)
}

func (s *ProviderService) CreateModelConfig(ctx context.Context, req *CreateModelConfigRequest) (*domain.ModelConfig, error) {
	modelConfig := &domain.ModelConfig{
		ProviderID:  req.ProviderID,
		Name:        req.Name,
		ModelName:   req.ModelName,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		ExtraConfig: req.ExtraConfig,
		IsEnabled:   true,
	}

	if err := s.modelConfigRepo.Create(modelConfig); err != nil {
		return nil, err
	}
	return modelConfig, nil
}

func (s *ProviderService) ListModelConfigs(ctx context.Context) ([]*domain.ModelConfig, error) {
	return s.modelConfigRepo.List()
}

func (s *ProviderService) GetProviderWithModelConfigs(ctx context.Context, id uint64) (*domain.Provider, []*domain.ModelConfig, error) {
	provider, err := s.providerRepo.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	modelConfigs, err := s.modelConfigRepo.ListByProviderID(id)
	if err != nil {
		return provider, nil, err
	}
	return provider, modelConfigs, nil
}
