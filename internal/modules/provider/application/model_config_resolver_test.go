package application

import (
	"context"
	"testing"

	"lattice-coding/internal/modules/provider/domain"
)

func TestModelConfigResolver_ResolveModelConfig_ShouldMapDomainToResolvedConfig(t *testing.T) {
	providerRepo := &resolverProviderRepo{
		providers: map[uint64]*domain.Provider{
			7: {
				ID:               7,
				ProviderType:     domain.ProviderTypeOpenAICompatible,
				BaseURL:          "https://api.example.com/v1",
				APIKeyCiphertext: "sk-test",
				Enabled:          true,
				Config:           `{"region":"local"}`,
			},
		},
	}
	modelConfigRepo := &resolverModelConfigRepo{
		configs: map[uint64]*domain.ModelConfig{
			11: {
				ID:         11,
				ProviderID: 7,
				Model:      "gpt-test",
				Params:     `{"temperature":0.2,"max_tokens":128}`,
				Enabled:    true,
			},
		},
	}

	resolver := NewModelConfigResolver(providerRepo, modelConfigRepo, nil)
	resolved, err := resolver.ResolveModelConfig(context.Background(), "11")
	if err != nil {
		t.Fatalf("ResolveModelConfig returned error: %v", err)
	}

	if resolved.ProviderID != "7" {
		t.Fatalf("ProviderID = %q, want 7", resolved.ProviderID)
	}
	if resolved.ProviderType != "openai_compatible" {
		t.Fatalf("ProviderType = %q, want openai_compatible", resolved.ProviderType)
	}
	if resolved.ModelConfigID != "11" {
		t.Fatalf("ModelConfigID = %q, want 11", resolved.ModelConfigID)
	}
	if resolved.ModelName != "gpt-test" {
		t.Fatalf("ModelName = %q, want gpt-test", resolved.ModelName)
	}
	if resolved.BaseURL != "https://api.example.com/v1" {
		t.Fatalf("BaseURL = %q", resolved.BaseURL)
	}
	if resolved.APIKey != "sk-test" {
		t.Fatalf("APIKey = %q, want sk-test", resolved.APIKey)
	}
	if resolved.Temperature != 0.2 {
		t.Fatalf("Temperature = %v, want 0.2", resolved.Temperature)
	}
	if resolved.MaxTokens != 128 {
		t.Fatalf("MaxTokens = %v, want 128", resolved.MaxTokens)
	}
	if resolved.Extra["temperature"].(float64) != 0.2 {
		t.Fatalf("temperature extra = %#v", resolved.Extra["temperature"])
	}
}

type resolverProviderRepo struct {
	providers map[uint64]*domain.Provider
}

func (r *resolverProviderRepo) Create(ctx context.Context, provider *domain.Provider) error {
	r.providers[provider.ID] = provider
	return nil
}

func (r *resolverProviderRepo) Update(ctx context.Context, provider *domain.Provider) error {
	r.providers[provider.ID] = provider
	return nil
}

func (r *resolverProviderRepo) FindByID(ctx context.Context, id uint64) (*domain.Provider, error) {
	return r.providers[id], nil
}

func (r *resolverProviderRepo) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.Provider], error) {
	items := make([]*domain.Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		items = append(items, provider)
	}
	return &domain.PageResult[*domain.Provider]{Items: items, Total: int64(len(items))}, nil
}

func (r *resolverProviderRepo) DeleteByID(ctx context.Context, id uint64) error {
	delete(r.providers, id)
	return nil
}

func (r *resolverProviderRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (r *resolverProviderRepo) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	r.providers[id].Enabled = enabled
	return nil
}

func (r *resolverProviderRepo) UpdateHealthStatus(ctx context.Context, id uint64, status domain.HealthStatus, lastError string) error {
	return nil
}

type resolverModelConfigRepo struct {
	configs map[uint64]*domain.ModelConfig
}

func (r *resolverModelConfigRepo) Create(ctx context.Context, modelConfig *domain.ModelConfig) error {
	r.configs[modelConfig.ID] = modelConfig
	return nil
}

func (r *resolverModelConfigRepo) Update(ctx context.Context, modelConfig *domain.ModelConfig) error {
	r.configs[modelConfig.ID] = modelConfig
	return nil
}

func (r *resolverModelConfigRepo) FindByID(ctx context.Context, id uint64) (*domain.ModelConfig, error) {
	return r.configs[id], nil
}

func (r *resolverModelConfigRepo) FindByProviderID(ctx context.Context, providerID uint64) ([]*domain.ModelConfig, error) {
	var result []*domain.ModelConfig
	for _, cfg := range r.configs {
		if cfg.ProviderID == providerID {
			result = append(result, cfg)
		}
	}
	return result, nil
}

func (r *resolverModelConfigRepo) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.ModelConfig], error) {
	items := make([]*domain.ModelConfig, 0, len(r.configs))
	for _, cfg := range r.configs {
		items = append(items, cfg)
	}
	return &domain.PageResult[*domain.ModelConfig]{Items: items, Total: int64(len(items))}, nil
}

func (r *resolverModelConfigRepo) DeleteByID(ctx context.Context, id uint64) error {
	delete(r.configs, id)
	return nil
}

func (r *resolverModelConfigRepo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (r *resolverModelConfigRepo) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	r.configs[id].Enabled = enabled
	return nil
}
