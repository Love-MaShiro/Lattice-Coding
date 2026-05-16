package persistence

import (
	"context"

	"lattice-coding/internal/modules/audit/domain"

	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepositoryImpl(db *gorm.DB) domain.AuditLogRepository {
	return &AuditLogRepositoryImpl{db: db}
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, log *domain.AuditLog) error {
	po := &AuditLogPO{
		RunID:        log.RunID,
		TraceID:      log.TraceID,
		EventType:    log.EventType,
		ToolName:     log.ToolName,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Message:      log.Message,
		PayloadJSON:  log.PayloadJSON,
	}
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	log.ID = po.ID
	log.CreatedAt = po.CreatedAt
	log.UpdatedAt = po.UpdatedAt
	return nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&AuditLogPO{})
}
