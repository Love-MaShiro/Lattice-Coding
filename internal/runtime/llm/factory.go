package llm

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"lattice-coding/internal/modules/provider/domain"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type LLMFactory struct {
	providerRepo       domain.ProviderRepository
	providerHealthRepo domain.ProviderHealthRepository
	modelConfigRepo    domain.ModelConfigRepository
}

func NewLLMFactory(
	providerRepo domain.ProviderRepository,
	providerHealthRepo domain.ProviderHealthRepository,
	modelConfigRepo domain.ModelConfigRepository,
) *LLMFactory {
	return &LLMFactory{
		providerRepo:       providerRepo,
		providerHealthRepo: providerHealthRepo,
		modelConfigRepo:    modelConfigRepo,
	}
}

func (f *LLMFactory) GetModelConfig(ctx context.Context, modelConfigID uint64) (*domain.ModelConfig, error) {
	return f.modelConfigRepo.FindByID(ctx, modelConfigID)
}

func (f *LLMFactory) ListModelConfigs(ctx context.Context) ([]*domain.ModelConfig, error) {
	result, err := f.modelConfigRepo.FindPage(ctx, &domain.PageRequest{Page: 1, PageSize: 100})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (f *LLMFactory) CreateChatModel(ctx context.Context, modelConfigID uint64) (model.ChatModel, error) {
	modelConfig, err := f.modelConfigRepo.FindByID(ctx, modelConfigID)
	if err != nil {
		return nil, err
	}

	provider, err := f.providerRepo.FindByID(ctx, modelConfig.ProviderID)
	if err != nil {
		return nil, err
	}

	if !provider.Enabled {
		return nil, ErrProviderDisabled
	}

	if !modelConfig.Enabled {
		return nil, ErrModelConfigDisabled
	}

	authConfig := decodeAuthConfig(provider.AuthConfigCiphertext)

	var apiKey string
	if authConfig != nil {
		apiKey = authConfig.APIKey
	}

	switch provider.ProviderType {
	case domain.ProviderTypeOpenAI, domain.ProviderTypeOpenAICompatible:
		return NewOpenAIChatModel(provider.BaseURL, apiKey, modelConfig.Model)
	case domain.ProviderTypeOllama:
		return NewOllamaChatModel(provider.BaseURL, modelConfig.Model)
	case domain.ProviderTypeClaude:
		return nil, ErrUnsupportedProviderType
	default:
		return nil, ErrUnsupportedProviderType
	}
}

func decodeAuthConfig(ciphertext string) *domain.AuthConfigData {
	if ciphertext == "" {
		return nil
	}
	var config domain.AuthConfigData
	if err := json.Unmarshal([]byte(ciphertext), &config); err != nil {
		return nil
	}
	return &config
}

func (f *LLMFactory) TestModel(ctx context.Context, modelConfigID uint64) (*HealthCheckResult, error) {
	return f.testModelWithTimeout(ctx, modelConfigID, 10*time.Second)
}

func (f *LLMFactory) testModelWithTimeout(ctx context.Context, modelConfigID uint64, timeout time.Duration) (*HealthCheckResult, error) {
	modelConfig, err := f.modelConfigRepo.FindByID(ctx, modelConfigID)
	if err != nil {
		return nil, err
	}

	provider, err := f.providerRepo.FindByID(ctx, modelConfig.ProviderID)
	if err != nil {
		return nil, err
	}

	result := &HealthCheckResult{
		ProviderID:    provider.ID,
		ModelConfigID: modelConfigID,
		CheckedAt:     time.Now(),
	}

	chatModel, err := f.CreateChatModel(ctx, modelConfigID)
	if err != nil {
		result.HealthStatus = string(domain.HealthStatusUnhealthy)
		result.ErrorMessage = summarizeError(err)
		return result, nil
	}

	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	startTime := time.Now()
	_, err = chatModel.Generate(testCtx, []*schema.Message{
		{Role: schema.User, Content: "请只回复 OK"},
	})
	latency := time.Since(startTime).Milliseconds()
	result.LatencyMs = latency

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			result.HealthStatus = string(domain.HealthStatusUnhealthy)
			result.ErrorCode = "TIMEOUT"
			result.ErrorMessage = "请求超时"
		} else {
			result.HealthStatus = string(domain.HealthStatusUnhealthy)
			result.ErrorCode = "CALL_ERROR"
			result.ErrorMessage = summarizeError(err)
		}
		return result, nil
	}

	result.Success = true
	result.HealthStatus = string(domain.HealthStatusHealthy)
	return result, nil
}

func (f *LLMFactory) TestProvider(ctx context.Context, providerID uint64) (*HealthCheckResult, error) {
	return f.testProviderWithTimeout(ctx, providerID, 10*time.Second)
}

func (f *LLMFactory) testProviderWithTimeout(ctx context.Context, providerID uint64, timeout time.Duration) (*HealthCheckResult, error) {
	modelConfigs, err := f.modelConfigRepo.FindByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}
	if len(modelConfigs) == 0 {
		return nil, ErrNoModelConfigFound
	}

	result := &HealthCheckResult{
		ProviderID: providerID,
		CheckedAt:  time.Now(),
	}

	var lastErr error
	for _, mc := range modelConfigs {
		if mc.Enabled {
			modelResult, err := f.testModelWithTimeout(ctx, mc.ID, timeout)
			if err != nil {
				lastErr = err
				continue
			}
			modelResult.ProviderID = providerID
			return modelResult, nil
		}
	}

	result.HealthStatus = string(domain.HealthStatusUnhealthy)
	result.ErrorCode = "NO_AVAILABLE_MODEL"
	result.ErrorMessage = "无可用的模型配置"
	if lastErr != nil {
		result.ErrorMessage = summarizeError(lastErr)
	}
	return result, nil
}

func summarizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()

	if strings.Contains(msg, "context deadline exceeded") || strings.Contains(msg, "Timeout") {
		return "请求超时"
	}
	if strings.Contains(msg, "connection refused") || strings.Contains(msg, "connect: connection refused") {
		return "连接被拒绝"
	}
	if strings.Contains(msg, "authentication") || strings.Contains(msg, "unauthorized") || strings.Contains(msg, "401") {
		return "认证失败"
	}
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "429") {
		return "请求频率超限"
	}
	if strings.Contains(msg, "500") || strings.Contains(msg, "502") || strings.Contains(msg, "503") {
		return "服务端错误"
	}

	if len(msg) > 100 {
		return msg[:100] + "..."
	}
	return msg
}
