package llm

import (
	"context"
	"errors"

	"lattice-coding/internal/modules/provider/domain"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type LLMFactory struct {
	providerRepo    domain.ProviderRepository
	modelConfigRepo domain.ModelConfigRepository
	cryptoService   domain.CryptoService
}

func NewLLMFactory(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	cryptoService domain.CryptoService,
) *LLMFactory {
	return &LLMFactory{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
		cryptoService:   cryptoService,
	}
}

func (f *LLMFactory) GetModelConfig(ctx context.Context, modelConfigID uint64) (*domain.ModelConfig, error) {
	return f.modelConfigRepo.GetByID(modelConfigID)
}

func (f *LLMFactory) ListModelConfigs(ctx context.Context) ([]*domain.ModelConfig, error) {
	return f.modelConfigRepo.List()
}

func (f *LLMFactory) CreateChatModel(ctx context.Context, modelConfigID uint64) (model.BaseChatModel, error) {
	modelConfig, err := f.modelConfigRepo.GetByID(modelConfigID)
	if err != nil {
		return nil, err
	}

	provider, err := f.providerRepo.GetByID(modelConfig.ProviderID)
	if err != nil {
		return nil, err
	}

	if !provider.IsEnabled || !modelConfig.IsEnabled {
		return nil, errors.New("provider or model config is disabled")
	}

	apiKey, err := f.cryptoService.Decrypt(provider.APIKeyCiphertext)
	if err != nil {
		return nil, err
	}

	switch provider.ProviderType {
	case domain.ProviderTypeOpenAI, domain.ProviderTypeOpenAICompatible:
		return NewOpenAIChatModel(provider.BaseURL, apiKey, modelConfig.ModelName)
	case domain.ProviderTypeOllama:
		return NewOllamaChatModel(provider.BaseURL, modelConfig.ModelName)
	case domain.ProviderTypeClaude:
		return nil, errors.New("claude provider not implemented yet")
	default:
		return nil, errors.New("unknown provider type")
	}
}

func (f *LLMFactory) TestModel(ctx context.Context, modelConfigID uint64) (bool, error) {
	chatModel, err := f.CreateChatModel(ctx, modelConfigID)
	if err != nil {
		return false, err
	}

	_, err = chatModel.Generate(ctx, []*schema.Message{
		{Role: schema.User, Content: "ping"},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *LLMFactory) TestProvider(ctx context.Context, providerID uint64) (bool, error) {
	modelConfigs, err := f.modelConfigRepo.ListByProviderID(providerID)
	if err != nil {
		return false, err
	}
	if len(modelConfigs) == 0 {
		return false, errors.New("no model config found")
	}
	return f.TestModel(ctx, modelConfigs[0].ID)
}
