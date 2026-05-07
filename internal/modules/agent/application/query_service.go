package application

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/agent/domain"
)

type QueryService struct {
	agentRepo     domain.AgentRepository
	agentToolRepo domain.AgentToolRepository
}

func NewQueryService(agentRepo domain.AgentRepository, agentToolRepo domain.AgentToolRepository) *QueryService {
	return &QueryService{
		agentRepo:     agentRepo,
		agentToolRepo: agentToolRepo,
	}
}

func (s *QueryService) GetAgent(ctx context.Context, id uint64) (*AgentDTO, error) {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return nil, errors.NotFoundErr("Agent 不存在")
	}

	return ToAgentDTO(agent), nil
}

func (s *QueryService) GetAgentDetail(ctx context.Context, id uint64) (*AgentDetailDTO, error) {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 失败")
	}
	if agent == nil {
		return nil, errors.NotFoundErr("Agent 不存在")
	}

	tools, err := s.agentToolRepo.FindByAgentID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 工具失败")
	}

	return ToAgentDetailDTO(agent, tools), nil
}

func (s *QueryService) ListAgents(ctx context.Context, query *AgentPageQuery) (*domain.PageResult[*AgentDTO], error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.agentRepo.FindPage(ctx, &domain.PageRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 列表失败")
	}

	agentIDs := make([]uint64, len(result.Items))
	for i, agent := range result.Items {
		agentIDs[i] = agent.ID
	}

	toolCounts, err := s.agentToolRepo.CountByAgentIDs(ctx, agentIDs)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Agent 工具数量失败")
	}

	dtos := make([]*AgentDTO, len(result.Items))
	for i, agent := range result.Items {
		toolCount := toolCounts[agent.ID]
		dtos[i] = ToAgentDTOWithToolCount(agent, toolCount)
	}

	return &domain.PageResult[*AgentDTO]{
		Items:    dtos,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
