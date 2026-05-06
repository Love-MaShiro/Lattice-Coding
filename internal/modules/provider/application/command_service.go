package application

import (
	"context"
	"lattice-coding/internal/common/crypto"
	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
)

type AgentReferenceChecker interface {
	HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error)
}

type CommandService struct {
	providerRepo      domain.ProviderRepository
	modelConfigRepo   domain.ModelConfigRepository
	encryptor         crypto.Encryptor
	agentRefChecker   AgentReferenceChecker
}

func NewCommandService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	encryptor crypto.Encryptor,
	agentRefChecker AgentReferenceChecker,
) *CommandService {
	return &CommandService{
		providerRepo:    providerRepo,
		modelConfigRepo: modelConfigRepo,
		encryptor:       encryptor,
		agentRefChecker: agentRefChecker,
	}
}

func (s *CommandService) CreateProvider(ctx context.Context, cmd *CreateProviderCommand) (*ProviderDTO, error) {
	exists, err := s.providerRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "检查 Provider 名称失败")
	}
	if exists {
		return nil, errors.AlreadyExistsErr("Provider 名称已存在")
	}

	apiKeyCiphertext := ""
	if cmd.APIKey != "" {
		ciphertext, err := s.encryptor.Encrypt(cmd.APIKey)
		if err != nil {
			return nil, errors.InternalWithErr(err, "加密 API Key 失败")
		}
		apiKeyCiphertext = ciphertext
	}

	authConfigCiphertext := ""
	if cmd.AuthConfig != "" {
		ciphertext, err := s.encryptor.Encrypt(cmd.AuthConfig)
		if err != nil {
			return nil, errors.InternalWithErr(err, "加密 Auth Config 失败")
		}
		authConfigCiphertext = ciphertext
	}

	provider := &domain.Provider{
		Name:                 cmd.Name,
		ProviderType:         domain.ProviderType(cmd.ProviderType),
		BaseURL:              cmd.BaseURL,
		AuthType:             domain.AuthType(cmd.AuthType),
		APIKeyCiphertext:     apiKeyCiphertext,
		AuthConfigCiphertext: authConfigCiphertext,
		Config:               cmd.Config,
		Enabled:              cmd.Enabled,
		HealthStatus:         domain.HealthStatusUnknown,
	}

	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "创建 Provider 失败")
	}

	return ToProviderDTO(provider), nil
}

func (s *CommandService) UpdateProvider(ctx context.Context, id uint64, cmd *UpdateProviderCommand) (*ProviderDTO, error) {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	if cmd.Name != "" && cmd.Name != provider.Name {
		exists, err := s.providerRepo.ExistsByName(ctx, cmd.Name)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "检查 Provider 名称失败")
		}
		if exists {
			return nil, errors.AlreadyExistsErr("Provider 名称已存在")
		}
		provider.Name = cmd.Name
	}

	if cmd.ProviderType != "" {
		provider.ProviderType = domain.ProviderType(cmd.ProviderType)
	}
	if cmd.BaseURL != "" {
		provider.BaseURL = cmd.BaseURL
	}
	if cmd.AuthType != "" {
		provider.AuthType = domain.AuthType(cmd.AuthType)
	}
	if cmd.Config != "" {
		provider.Config = cmd.Config
	}
	if cmd.Enabled != nil {
		provider.Enabled = *cmd.Enabled
	}

	if cmd.APIKey != "" {
		ciphertext, err := s.encryptor.Encrypt(cmd.APIKey)
		if err != nil {
			return nil, errors.InternalWithErr(err, "加密 API Key 失败")
		}
		provider.APIKeyCiphertext = ciphertext
	}

	if cmd.AuthConfig != "" {
		ciphertext, err := s.encryptor.Encrypt(cmd.AuthConfig)
		if err != nil {
			return nil, errors.InternalWithErr(err, "加密 Auth Config 失败")
		}
		provider.AuthConfigCiphertext = ciphertext
	}

	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "更新 Provider 失败")
	}

	return ToProviderDTO(provider), nil
}

func (s *CommandService) DeleteProvider(ctx context.Context, id uint64) error {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return errors.NotFoundErr("Provider 不存在")
	}

	modelConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询关联模型配置失败")
	}
	if len(modelConfigs) > 0 {
		return errors.ForbiddenErr("该 Provider 下仍有关联的模型配置，请先删除所有模型配置")
	}

	if err := s.providerRepo.DeleteByID(ctx, id); err != nil {
		return errors.DatabaseErrWithErr(err, "删除 Provider 失败")
	}

	return nil
}

func (s *CommandService) EnableProvider(ctx context.Context, id uint64) error {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return errors.NotFoundErr("Provider 不存在")
	}

	if err := s.providerRepo.UpdateEnabled(ctx, id, true); err != nil {
		return errors.DatabaseErrWithErr(err, "启用 Provider 失败")
	}

	return nil
}

func (s *CommandService) DisableProvider(ctx context.Context, id uint64) error {
	provider, err := s.providerRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return errors.NotFoundErr("Provider 不存在")
	}

	if err := s.providerRepo.UpdateEnabled(ctx, id, false); err != nil {
		return errors.DatabaseErrWithErr(err, "禁用 Provider 失败")
	}

	return nil
}

func (s *CommandService) CreateModelConfig(ctx context.Context, cmd *CreateModelConfigCommand) (*ModelConfigDTO, error) {
	provider, err := s.providerRepo.FindByID(ctx, cmd.ProviderID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	modelConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, cmd.ProviderID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
	}

	for _, mc := range modelConfigs {
		if mc.Name == cmd.Name {
			return nil, errors.AlreadyExistsErr("同一 Provider 下模型配置名称已存在")
		}
	}

	modelConfig := &domain.ModelConfig{
		ProviderID:   cmd.ProviderID,
		Name:         cmd.Name,
		Model:        cmd.Model,
		ModelType:    domain.ModelType(cmd.ModelType),
		Params:       cmd.Params,
		Capabilities: cmd.Capabilities,
		IsDefault:    cmd.IsDefault,
		Enabled:      cmd.Enabled,
	}

	if cmd.IsDefault {
		for _, mc := range modelConfigs {
			if mc.IsDefault && mc.ModelType == modelConfig.ModelType {
				return nil, errors.AlreadyExistsErr("同一 Provider 下相同模型类型已有默认模型")
			}
		}
	}

	if err := s.modelConfigRepo.Create(ctx, modelConfig); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "创建模型配置失败")
	}

	return ToModelConfigDTO(modelConfig), nil
}

func (s *CommandService) UpdateModelConfig(ctx context.Context, id uint64, cmd *UpdateModelConfigCommand) (*ModelConfigDTO, error) {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return nil, errors.NotFoundErr("模型配置不存在")
	}

	if cmd.Name != "" && cmd.Name != modelConfig.Name {
		providerConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, modelConfig.ProviderID)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
		}
		for _, mc := range providerConfigs {
			if mc.ID != id && mc.Name == cmd.Name {
				return nil, errors.AlreadyExistsErr("同一 Provider 下模型配置名称已存在")
			}
		}
		modelConfig.Name = cmd.Name
	}

	if cmd.Model != "" {
		modelConfig.Model = cmd.Model
	}
	if cmd.ModelType != "" {
		modelConfig.ModelType = domain.ModelType(cmd.ModelType)
	}
	if cmd.Params != "" {
		modelConfig.Params = cmd.Params
	}
	if cmd.Capabilities != "" {
		modelConfig.Capabilities = cmd.Capabilities
	}
	if cmd.Enabled != nil {
		modelConfig.Enabled = *cmd.Enabled
	}

	if cmd.IsDefault != nil && *cmd.IsDefault {
		providerConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, modelConfig.ProviderID)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
		}
		for _, mc := range providerConfigs {
			if mc.ID != id && mc.IsDefault && mc.ModelType == modelConfig.ModelType {
				mc.IsDefault = false
				if err := s.modelConfigRepo.Update(ctx, mc); err != nil {
					return nil, errors.DatabaseErrWithErr(err, "更新默认模型配置失败")
				}
			}
		}
		modelConfig.IsDefault = true
	}

	if err := s.modelConfigRepo.Update(ctx, modelConfig); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "更新模型配置失败")
	}

	return ToModelConfigDTO(modelConfig), nil
}

func (s *CommandService) DeleteModelConfig(ctx context.Context, id uint64) error {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return errors.NotFoundErr("模型配置不存在")
	}

	if s.agentRefChecker != nil {
		hasRef, err := s.agentRefChecker.HasModelConfigReferences(ctx, id)
		if err != nil {
			return errors.InternalWithErr(err, "检查 Agent 引用失败")
		}
		if hasRef {
			return errors.ForbiddenErr("该模型配置正在被 Agent 引用，无法删除")
		}
	}

	if err := s.modelConfigRepo.DeleteByID(ctx, id); err != nil {
		return errors.DatabaseErrWithErr(err, "删除模型配置失败")
	}

	return nil
}

func (s *CommandService) EnableModelConfig(ctx context.Context, id uint64) error {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return errors.NotFoundErr("模型配置不存在")
	}

	if err := s.modelConfigRepo.UpdateEnabled(ctx, id, true); err != nil {
		return errors.DatabaseErrWithErr(err, "启用模型配置失败")
	}

	return nil
}

func (s *CommandService) DisableModelConfig(ctx context.Context, id uint64) error {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return errors.NotFoundErr("模型配置不存在")
	}

	if err := s.modelConfigRepo.UpdateEnabled(ctx, id, false); err != nil {
		return errors.DatabaseErrWithErr(err, "禁用模型配置失败")
	}

	return nil
}

func (s *CommandService) SetDefaultModelConfig(ctx context.Context, id uint64) error {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return errors.NotFoundErr("模型配置不存在")
	}

	providerConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, modelConfig.ProviderID)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询模型配置列表失败")
	}

	for _, mc := range providerConfigs {
		if mc.ID != id && mc.IsDefault && mc.ModelType == modelConfig.ModelType {
			mc.IsDefault = false
			if err := s.modelConfigRepo.Update(ctx, mc); err != nil {
				return errors.DatabaseErrWithErr(err, "更新默认模型配置失败")
			}
		}
	}

	modelConfig.IsDefault = true
	if err := s.modelConfigRepo.Update(ctx, modelConfig); err != nil {
		return errors.DatabaseErrWithErr(err, "设置默认模型配置失败")
	}

	return nil
}
