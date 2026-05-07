package application

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/agent/domain"
)

type ModelConfigChecker interface {
	CheckModelConfigEnabled(ctx context.Context, modelConfigID uint64) error
}

type CommandService struct {
	agentRepo          domain.AgentRepository
	modelConfigChecker ModelConfigChecker
}

func NewCommandService(agentRepo domain.AgentRepository, modelConfigChecker ModelConfigChecker) *CommandService {
	return &CommandService{
		agentRepo:          agentRepo,
		modelConfigChecker: modelConfigChecker,
	}
}

func (s *CommandService) CreateAgent(ctx context.Context, cmd *CreateAgentCommand) (*AgentDTO, error) {
	if cmd.ModelConfigID == 0 {
		return nil, errors.InvalidArg("model_config_id 不能为空")
	}

	if err := s.modelConfigChecker.CheckModelConfigEnabled(ctx, cmd.ModelConfigID); err != nil {
		return nil, err
	}

	exists, err := s.agentRepo.ExistsByName(ctx, cmd.Name, 0)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "检查 Agent 名称失败")
	}
	if exists {
		return nil, errors.AlreadyExistsErr("Agent 名称已存在")
	}

	agent := &domain.Agent{
		Name:            cmd.Name,
		Description:     cmd.Description,
		AgentType:       domain.AgentType(cmd.AgentType),
		ModelConfigID:   cmd.ModelConfigID,
		SystemPrompt:    cmd.SystemPrompt,
		Temperature:     cmd.Temperature,
		TopP:            cmd.TopP,
		MaxTokens:       cmd.MaxTokens,
		MaxContextTurns: cmd.MaxContextTurns,
		MaxSteps:        cmd.MaxSteps,
		Enabled:         cmd.Enabled,
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
		exists, err := s.agentRepo.ExistsByName(ctx, cmd.Name, id)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "检查 Agent 名称失败")
		}
		if exists {
			return nil, errors.AlreadyExistsErr("Agent 名称已存在")
		}
		agent.Name = cmd.Name
	}

	if cmd.Description != "" {
		agent.Description = cmd.Description
	}
	if cmd.AgentType != "" {
		agent.AgentType = domain.AgentType(cmd.AgentType)
	}
	if cmd.ModelConfigID > 0 && cmd.ModelConfigID != agent.ModelConfigID {
		if err := s.modelConfigChecker.CheckModelConfigEnabled(ctx, cmd.ModelConfigID); err != nil {
			return nil, err
		}
		agent.ModelConfigID = cmd.ModelConfigID
	}
	if cmd.SystemPrompt != "" {
		agent.SystemPrompt = cmd.SystemPrompt
	}
	if cmd.Temperature > 0 {
		agent.Temperature = cmd.Temperature
	}
	if cmd.TopP > 0 {
		agent.TopP = cmd.TopP
	}
	if cmd.MaxTokens > 0 {
		agent.MaxTokens = cmd.MaxTokens
	}
	if cmd.MaxContextTurns > 0 {
		agent.MaxContextTurns = cmd.MaxContextTurns
	}
	if cmd.MaxSteps > 0 {
		agent.MaxSteps = cmd.MaxSteps
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
