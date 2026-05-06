package persistence

import (
	"context"

	"lattice-coding/internal/modules/agent/domain"

	"gorm.io/gorm"
)

type AgentRepositoryImpl struct {
	db *gorm.DB
}

func NewAgentRepositoryImpl(db *gorm.DB) domain.AgentRepository {
	return &AgentRepositoryImpl{db: db}
}

func (r *AgentRepositoryImpl) Create(ctx context.Context, agent *domain.Agent) error {
	po := &AgentPO{}
	ConvertAgentToPO(agent, po)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *AgentRepositoryImpl) Update(ctx context.Context, agent *domain.Agent) error {
	po := &AgentPO{}
	ConvertAgentToPO(agent, po)
	return r.db.WithContext(ctx).Model(po).Omit("created_at").Updates(po).Error
}

func (r *AgentRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.Agent, error) {
	var po AgentPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return ConvertPOToAgent(&po), nil
}

func (r *AgentRepositoryImpl) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.Agent], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&AgentPO{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	var pos []AgentPO
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&pos).Error; err != nil {
		return nil, err
	}

	items := make([]*domain.Agent, len(pos))
	for i := range pos {
		items[i] = ConvertPOToAgent(&pos[i])
	}

	return &domain.PageResult[*domain.Agent]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *AgentRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&AgentPO{}, id).Error
}

func (r *AgentRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&AgentPO{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AgentRepositoryImpl) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	return r.db.WithContext(ctx).Model(&AgentPO{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&AgentPO{})
}
