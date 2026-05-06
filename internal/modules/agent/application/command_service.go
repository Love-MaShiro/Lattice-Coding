package application

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/agent/domain"
)

type ModelConfigGetter interface {
	GetModelConfig(ctx context.Context, id uint64) (uint64, error)
	GetProvider(ctx context.Context, id uint64) error
}

type CommandService struct {
	agentRepo         domain.AgentRepository
	modelConfigGetter ModelConfigGetter
}

func NewCommandService(
	agentRepo domain.AgentRepository,
	modelConfigGetter ModelConfigGetter,
) *CommandService {
	return &CommandService{
		agentRepo:         agentRepo,
		modelConfigGetter: modelConfigGetter,
	}
}

func (s *CommandService) CreateAgent(ctx context.Context, cmd *CreateAgentCommand) (*AgentDTO, error) {
	exists, err := s.agentRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "检查 Agent 名称失败")
	}
	if exists {
		return nil, errors.AlreadyExistsErr("Agent 名称已存在")
	}

	providerID, err := s.modelConfigGetter.GetModelConfig(ctx, cmd.ModelConfigID)
	if err != nil {
		return nil, err
	}
	if providerID != cmd.ProviderID {
		return nil, errors.ForbiddenErr("ModelConfig 不属于指定的 Provider")
	}

	if err := s.modelConfigGetter.GetProvider(ctx, cmd.ProviderID); err != nil {
		return nil, err
	}

	agent := &domain.Agent{
		Name:          cmd.Name,
		Description:   cmd.Description,
		ProviderID:    cmd.ProviderID,
		ModelConfigID: cmd.ModelConfigID,
		SystemPrompt:  cmd.SystemPrompt,
		Tools:         cmd.Tools,
		MaxSteps:      cmd.MaxSteps,
		Timeout:       cmd.Timeout,
		Enabled:       cmd.Enabled,
	}

	if err := s.agentRepo.Create(ctx, agent); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "创建 Agent 失败")
	}

	return ToAgentDTO(agent), nil
}

func (s *CommandService) UpdateAgent(ctx context.Context, id uint64, cmd *UpdateAgentCommand) (*AgentDTO, error) {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return nil, errors.NotFoundErr("Agent 不存在")
	}

	if cmd.Name != "" && cmd.Name != agent.Name {
		exists, err := s.agentRepo.ExistsByName(ctx, cmd.Name)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "检查 Agent 名称失败")
		}
		if exists {
			return nil, errors.AlreadyExistsErr("Agent 名称已存在")
		}
		agent.Name = cmd.Name
	}

	if cmd.ProviderID > 0 && cmd.ProviderID != agent.ProviderID {
		if err := s.modelConfigGetter.GetProvider(ctx, cmd.ProviderID); err != nil {
			return nil, err
		}
		agent.ProviderID = cmd.ProviderID
	}

	if cmd.ModelConfigID > 0 && cmd.ModelConfigID != agent.ModelConfigID {
		providerID, err := s.modelConfigGetter.GetModelConfig(ctx, cmd.ModelConfigID)
		if err != nil {
			return nil, err
		}
		effectiveProviderID := cmd.ProviderID
		if effectiveProviderID == 0 {
			effectiveProviderID = agent.ProviderID
		}
		if providerID != effectiveProviderID {
			return nil, errors.ForbiddenErr("ModelConfig 不属于指定的 Provider")
		}
		agent.ModelConfigID = cmd.ModelConfigID
	}

	if cmd.Description != "" {
		agent.Description = cmd.Description
	}
	if cmd.SystemPrompt != "" {
		agent.SystemPrompt = cmd.SystemPrompt
	}
	if cmd.Tools != "" {
		agent.Tools = cmd.Tools
	}
	if cmd.MaxSteps > 0 {
		agent.MaxSteps = cmd.MaxSteps
	}
	if cmd.Timeout > 0 {
		agent.Timeout = cmd.Timeout
	}
	if cmd.Enabled != nil {
		agent.Enabled = *cmd.Enabled
	}

	if err := s.agentRepo.Update(ctx, agent); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "更新 Agent 失败")
	}

	return ToAgentDTO(agent), nil
}

func (s *CommandService) DeleteAgent(ctx context.Context, id uint64) error {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return errors.NotFoundErr("Agent 不存在")
	}

	if err := s.agentRepo.DeleteByID(ctx, id); err != nil {
		return errors.DatabaseErrWithErr(err, "删除 Agent 失败")
	}

	return nil
}

func (s *CommandService) EnableAgent(ctx context.Context, id uint64) error {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return errors.NotFoundErr("Agent 不存在")
	}

	if err := s.agentRepo.UpdateEnabled(ctx, id, true); err != nil {
		return errors.DatabaseErrWithErr(err, "启用 Agent 失败")
	}

	return nil
}

func (s *CommandService) DisableAgent(ctx context.Context, id uint64) error {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return errors.NotFoundErr("Agent 不存在")
	}

	if err := s.agentRepo.UpdateEnabled(ctx, id, false); err != nil {
		return errors.DatabaseErrWithErr(err, "禁用 Agent 失败")
	}

	return nil
}
