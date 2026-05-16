package application

import (
	"context"

	"lattice-coding/internal/modules/run/domain"
)

type QueryService struct {
	runRepo        domain.RunRepository
	invocationRepo domain.ToolInvocationRepository
}

func NewQueryService(runRepo domain.RunRepository, invocationRepo domain.ToolInvocationRepository) *QueryService {
	return &QueryService{runRepo: runRepo, invocationRepo: invocationRepo}
}

func (s *QueryService) ListRuns(ctx context.Context, query RunPageQuery) (*PageResult[*RunDTO], error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	result, err := s.runRepo.FindPage(ctx, domain.PageRequest{Page: query.Page, PageSize: query.PageSize})
	if err != nil {
		return nil, err
	}
	return &PageResult[*RunDTO]{
		Items:    ToRunDTOs(result.Items),
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *QueryService) GetRun(ctx context.Context, id string) (*RunDTO, error) {
	run, err := s.runRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToRunDTO(run), nil
}

func (s *QueryService) ListToolInvocations(ctx context.Context, runID string) ([]*ToolInvocationDTO, error) {
	invocations, err := s.invocationRepo.FindByRunID(ctx, runID)
	if err != nil {
		return nil, err
	}
	return ToToolInvocationDTOs(invocations), nil
}
