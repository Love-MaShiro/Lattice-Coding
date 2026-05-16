package llm

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type LLMFactory struct {
	resolver ModelConfigResolver
}

func NewLLMFactory(resolver ModelConfigResolver) *LLMFactory {
	return &LLMFactory{resolver: resolver}
}

func (f *LLMFactory) ResolveModelConfig(ctx context.Context, modelConfigID uint64) (*ResolvedModelConfig, error) {
	if f == nil || f.resolver == nil {
		return nil, ErrNoProvider
	}
	return f.resolver.ResolveModelConfig(ctx, strconv.FormatUint(modelConfigID, 10))
}

func (f *LLMFactory) CreateChatModel(ctx context.Context, modelConfigID uint64) (model.ChatModel, error) {
	config, err := f.ResolveModelConfig(ctx, modelConfigID)
	if err != nil {
		return nil, err
	}
	return NewChatModelFromResolvedConfig(config)
}

func NewChatModelFromResolvedConfig(config *ResolvedModelConfig) (model.ChatModel, error) {
	if config == nil {
		return nil, ErrNoProvider
	}

	switch config.ProviderType {
	case "openai", "openai_compatible":
		return NewOpenAIChatModel(config.BaseURL, config.APIKey, config.ModelName)
	case "ollama":
		return NewOllamaChatModel(config.BaseURL, config.ModelName)
	case "claude":
		return nil, ErrUnsupportedProviderType
	default:
		return nil, ErrUnsupportedProviderType
	}
}

func (f *LLMFactory) TestModel(ctx context.Context, modelConfigID uint64) (*HealthCheckResult, error) {
	return f.testModelWithTimeout(ctx, modelConfigID, 10*time.Second)
}

func (f *LLMFactory) testModelWithTimeout(ctx context.Context, modelConfigID uint64, timeout time.Duration) (*HealthCheckResult, error) {
	config, err := f.ResolveModelConfig(ctx, modelConfigID)
	if err != nil {
		return nil, err
	}

	providerID, _ := strconv.ParseUint(config.ProviderID, 10, 64)
	resolvedModelConfigID, _ := strconv.ParseUint(config.ModelConfigID, 10, 64)
	if resolvedModelConfigID == 0 {
		resolvedModelConfigID = modelConfigID
	}

	result := &HealthCheckResult{
		ProviderID:    providerID,
		ModelConfigID: resolvedModelConfigID,
		CheckedAt:     time.Now(),
	}

	chatModel, err := NewChatModelFromResolvedConfig(config)
	if err != nil {
		result.HealthStatus = "unhealthy"
		result.ErrorMessage = summarizeError(err)
		return result, nil
	}

	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	startTime := time.Now()
	_, err = chatModel.Generate(testCtx, []*schema.Message{
		{Role: schema.User, Content: "Please reply OK"},
	})
	result.LatencyMs = time.Since(startTime).Milliseconds()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			result.HealthStatus = "unhealthy"
			result.ErrorCode = "TIMEOUT"
			result.ErrorMessage = "request timeout"
		} else {
			result.HealthStatus = "unhealthy"
			result.ErrorCode = "CALL_ERROR"
			result.ErrorMessage = summarizeError(err)
		}
		return result, nil
	}

	result.Success = true
	result.HealthStatus = "healthy"
	return result, nil
}

func (f *LLMFactory) TestProvider(ctx context.Context, providerID uint64) (*HealthCheckResult, error) {
	return f.testProviderWithTimeout(ctx, providerID, 10*time.Second)
}

func (f *LLMFactory) testProviderWithTimeout(ctx context.Context, providerID uint64, timeout time.Duration) (*HealthCheckResult, error) {
	if f == nil || f.resolver == nil {
		return nil, ErrNoProvider
	}

	providerResolver, ok := f.resolver.(ProviderDefaultModelResolver)
	if !ok {
		return nil, ErrNoModelConfigFound
	}

	config, err := providerResolver.ResolveProviderDefaultModel(ctx, strconv.FormatUint(providerID, 10))
	if err != nil {
		return nil, err
	}

	modelConfigID, _ := strconv.ParseUint(config.ModelConfigID, 10, 64)
	if modelConfigID == 0 {
		return nil, ErrNoModelConfigFound
	}
	return f.testModelWithTimeout(ctx, modelConfigID, timeout)
}

func summarizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()

	if strings.Contains(msg, "context deadline exceeded") || strings.Contains(msg, "Timeout") {
		return "request timeout"
	}
	if strings.Contains(msg, "connection refused") || strings.Contains(msg, "connect: connection refused") {
		return "connection refused"
	}
	if strings.Contains(msg, "authentication") || strings.Contains(msg, "unauthorized") || strings.Contains(msg, "401") {
		return "authentication failed"
	}
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "429") {
		return "rate limit exceeded"
	}
	if strings.Contains(msg, "500") || strings.Contains(msg, "502") || strings.Contains(msg, "503") {
		return "server error"
	}

	if len(msg) > 100 {
		return msg[:100] + "..."
	}
	return msg
}
