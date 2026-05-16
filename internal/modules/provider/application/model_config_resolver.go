package application

import (
	"context"
	"encoding/json"
	"strconv"

	"lattice-coding/internal/common/crypto"
	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
	"lattice-coding/internal/runtime/llm"
)

type ModelConfigResolver struct {
	providerRepo    domain.ProviderRepository
	modelConfigRepo domain.ModelConfigRepository
	encryptor       crypto.Encryptor
}

func NewModelConfigResolver(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	encryptor crypto.Encryptor,
) *ModelConfigResolver {
	if encryptor == nil {
		encryptor = crypto.NewNoopEncryptor()
	}
	return &ModelConfigResolver{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
		encryptor:       encryptor,
	}
}

func (r *ModelConfigResolver) ResolveModelConfig(ctx context.Context, modelConfigID string) (*llm.ResolvedModelConfig, error) {
	id, err := strconv.ParseUint(modelConfigID, 10, 64)
	if err != nil {
		return nil, errors.InvalidArg("invalid model config id")
	}
	modelConfig, err := r.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "query model config failed")
	}
	if modelConfig == nil {
		return nil, errors.NotFoundErr("model config not found")
	}
	return r.resolve(ctx, modelConfig)
}

func (r *ModelConfigResolver) ResolveAgentDefaultModel(ctx context.Context, agentID string) (*llm.ResolvedModelConfig, error) {
	return nil, errors.NotFoundErr("agent default model resolver is not connected")
}

func (r *ModelConfigResolver) ResolveProviderDefaultModel(ctx context.Context, providerID string) (*llm.ResolvedModelConfig, error) {
	id, err := strconv.ParseUint(providerID, 10, 64)
	if err != nil {
		return nil, errors.InvalidArg("invalid provider id")
	}
	modelConfigs, err := r.modelConfigRepo.FindByProviderID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "query provider model configs failed")
	}

	var selected *domain.ModelConfig
	for _, cfg := range modelConfigs {
		if cfg.Enabled && cfg.IsDefault {
			selected = cfg
			break
		}
		if selected == nil && cfg.Enabled {
			selected = cfg
		}
	}
	if selected == nil {
		return nil, errors.NotFoundErr("no enabled model config found")
	}
	return r.resolve(ctx, selected)
}

func (r *ModelConfigResolver) resolve(ctx context.Context, modelConfig *domain.ModelConfig) (*llm.ResolvedModelConfig, error) {
	provider, err := r.providerRepo.FindByID(ctx, modelConfig.ProviderID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "query provider failed")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("provider not found")
	}
	if !provider.Enabled {
		return nil, errors.ForbiddenErr("provider is disabled")
	}
	if !modelConfig.Enabled {
		return nil, errors.ForbiddenErr("model config is disabled")
	}

	apiKey, err := r.resolveAPIKey(provider)
	if err != nil {
		return nil, err
	}

	extra := map[string]interface{}{}
	if modelConfig.Params != "" {
		_ = json.Unmarshal([]byte(modelConfig.Params), &extra)
	}
	if provider.Config != "" {
		var providerConfig map[string]interface{}
		if err := json.Unmarshal([]byte(provider.Config), &providerConfig); err == nil {
			extra["provider_config"] = providerConfig
		}
	}

	return &llm.ResolvedModelConfig{
		ProviderID:    strconv.FormatUint(provider.ID, 10),
		ProviderType:  string(provider.ProviderType),
		ModelConfigID: strconv.FormatUint(modelConfig.ID, 10),
		ModelName:     modelConfig.Model,
		BaseURL:       provider.BaseURL,
		APIKey:        apiKey,
		Temperature:   floatFromExtra(extra, "temperature"),
		TopP:          floatFromExtra(extra, "top_p"),
		MaxTokens:     intFromExtra(extra, "max_tokens"),
		Extra:         extra,
	}, nil
}

func (r *ModelConfigResolver) resolveAPIKey(provider *domain.Provider) (string, error) {
	if provider.APIKeyCiphertext != "" {
		return r.encryptor.Decrypt(provider.APIKeyCiphertext)
	}
	if provider.AuthConfigCiphertext == "" {
		return "", nil
	}

	plaintext, err := r.encryptor.Decrypt(provider.AuthConfigCiphertext)
	if err != nil {
		return "", err
	}

	var authConfig domain.AuthConfigData
	if err := json.Unmarshal([]byte(plaintext), &authConfig); err != nil {
		return "", nil
	}
	return authConfig.APIKey, nil
}

func floatFromExtra(extra map[string]interface{}, key string) float64 {
	value, ok := extra[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case json.Number:
		result, _ := typed.Float64()
		return result
	default:
		return 0
	}
}

func intFromExtra(extra map[string]interface{}, key string) int {
	value, ok := extra[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case json.Number:
		result, _ := typed.Int64()
		return int(result)
	default:
		return 0
	}
}
