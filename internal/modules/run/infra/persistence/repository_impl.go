package persistence

import (
	"context"

	"lattice-coding/internal/modules/run/domain"

	"gorm.io/gorm"
)

type RunRepositoryImpl struct {
	db *gorm.DB
}

func NewRunRepositoryImpl(db *gorm.DB) domain.RunRepository {
	return &RunRepositoryImpl{db: db}
}

func (r *RunRepositoryImpl) Create(ctx context.Context, run *domain.Run) error {
	po := ToRunPO(run)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*run = *ToRunDomain(po)
	return nil
}

func (r *RunRepositoryImpl) Update(ctx context.Context, run *domain.Run) error {
	updates := map[string]interface{}{
		"agent_id":     run.AgentID,
		"session_id":   run.SessionID,
		"workflow_id":  run.WorkflowID,
		"status":       run.Status,
		"input":        run.Input,
		"output":       run.Output,
		"error":        run.Error,
		"started_at":   run.StartedAt,
		"completed_at": run.CompletedAt,
	}
	return r.db.WithContext(ctx).Model(&RunPO{}).Where("id = ?", run.ID).Updates(updates).Error
}

func (r *RunRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.Run, error) {
	var po RunPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ToRunDomain(&po), nil
}

func (r *RunRepositoryImpl) FindPage(ctx context.Context, req domain.PageRequest) (*domain.PageResult[*domain.Run], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&RunPO{}).Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (req.Page - 1) * req.PageSize
	var pos []RunPO
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Order("started_at DESC, created_at DESC").Find(&pos).Error; err != nil {
		return nil, err
	}
	items := make([]*domain.Run, 0, len(pos))
	for i := range pos {
		items = append(items, ToRunDomain(&pos[i]))
	}
	return &domain.PageResult[*domain.Run]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

type ToolInvocationRepositoryImpl struct {
	db *gorm.DB
}

func NewToolInvocationRepositoryImpl(db *gorm.DB) domain.ToolInvocationRepository {
	return &ToolInvocationRepositoryImpl{db: db}
}

func (r *ToolInvocationRepositoryImpl) Create(ctx context.Context, invocation *domain.ToolInvocation) error {
	po := ToInvocationPO(invocation)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*invocation = *ToInvocationDomain(po)
	return nil
}

func (r *ToolInvocationRepositoryImpl) Update(ctx context.Context, invocation *domain.ToolInvocation) error {
	updates := map[string]interface{}{
		"run_id":          invocation.RunID,
		"node_id":         invocation.NodeID,
		"tool_name":       invocation.ToolName,
		"output_json":     invocation.OutputJSON,
		"is_error":        invocation.IsError,
		"latency_ms":      invocation.LatencyMs,
		"status":          invocation.Status,
		"full_result_ref": invocation.FullResultRef,
		"completed_at":    invocation.CompletedAt,
	}
	return r.db.WithContext(ctx).Model(&ToolInvocationPO{}).Where("id = ?", invocation.ID).Updates(updates).Error
}

func (r *ToolInvocationRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.ToolInvocation, error) {
	var po ToolInvocationPO
	if err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return ToInvocationDomain(&po), nil
}

func (r *ToolInvocationRepositoryImpl) FindByRunID(ctx context.Context, runID string) ([]*domain.ToolInvocation, error) {
	var pos []ToolInvocationPO
	if err := r.db.WithContext(ctx).Where("run_id = ?", runID).Order("started_at ASC, created_at ASC").Find(&pos).Error; err != nil {
		return nil, err
	}
	items := make([]*domain.ToolInvocation, 0, len(pos))
	for i := range pos {
		items = append(items, ToInvocationDomain(&pos[i]))
	}
	return items, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&RunPO{}, &ToolInvocationPO{})
}
