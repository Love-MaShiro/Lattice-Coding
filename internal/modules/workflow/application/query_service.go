package application

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/workflow/domain"
)

type QueryService struct {
	repo domain.WorkflowRepository
}

func NewQueryService(repo domain.WorkflowRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) GetWorkflow(ctx context.Context, id uint64) (*WorkflowDTO, error) {
	if id == 0 {
		return nil, errors.InvalidArg("workflow id is required")
	}
	workflow, err := s.repo.FindByIDWithGraph(ctx, id)
	if err != nil {
		return nil, errors.NotFoundErr("workflow not found")
	}
	return ToWorkflowDTO(workflow), nil
}

func (s *QueryService) ListWorkflows(ctx context.Context, query WorkflowPageQuery) (*WorkflowPageDTO, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}
	result, err := s.repo.FindPage(ctx, &domain.PageRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
		Keyword:  query.Keyword,
		Status:   query.Status,
	})
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "list workflows failed")
	}
	items := make([]*WorkflowDTO, len(result.Items))
	for i := range result.Items {
		items[i] = ToWorkflowSummaryDTO(result.Items[i])
	}
	return &WorkflowPageDTO{
		Items:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
