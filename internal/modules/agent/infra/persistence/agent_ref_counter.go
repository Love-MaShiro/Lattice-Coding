package persistence

import (
	"context"

	"lattice-coding/internal/modules/agent/domain"

	"gorm.io/gorm"
)

type AgentRefCounter struct {
	db *gorm.DB
}

func NewAgentRefCounter(db *gorm.DB) domain.AgentReferenceChecker {
	return &AgentRefCounter{db: db}
}

func (r *AgentRefCounter) HasModelConfigReferences(ctx context.Context, modelConfigID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&AgentPO{}).
		Where("model_config_id = ? AND deleted_at IS NULL", modelConfigID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
