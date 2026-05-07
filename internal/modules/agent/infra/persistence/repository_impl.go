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

func (r *AgentRepositoryImpl) ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&AgentPO{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AgentRepositoryImpl) UpdateEnabled(ctx context.Context, id uint64, enabled bool) error {
	return r.db.WithContext(ctx).Model(&AgentPO{}).Where("id = ?", id).Update("enabled", enabled).Error
}

type AgentToolRepositoryImpl struct {
	db *gorm.DB
}

func NewAgentToolRepositoryImpl(db *gorm.DB) domain.AgentToolRepository {
	return &AgentToolRepositoryImpl{db: db}
}

func (r *AgentToolRepositoryImpl) BatchCreate(ctx context.Context, agentTools []*domain.AgentTool) error {
	if len(agentTools) == 0 {
		return nil
	}
	pos := make([]AgentToolPO, len(agentTools))
	for i, t := range agentTools {
		ConvertAgentToolToPO(t, &pos[i])
	}
	return r.db.WithContext(ctx).Create(&pos).Error
}

func (r *AgentToolRepositoryImpl) DeleteByAgentID(ctx context.Context, agentID uint64) error {
	return r.db.WithContext(ctx).Where("agent_id = ?", agentID).Delete(&AgentToolPO{}).Error
}

func (r *AgentToolRepositoryImpl) FindByAgentID(ctx context.Context, agentID uint64) ([]*domain.AgentTool, error) {
	var pos []AgentToolPO
	if err := r.db.WithContext(ctx).Where("agent_id = ?", agentID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return ConvertPOsToAgentTools(pos), nil
}

func (r *AgentToolRepositoryImpl) CountByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error) {
	if len(agentIDs) == 0 {
		return map[uint64]int64{}, nil
	}

	type countResult struct {
		AgentID uint64 `gorm:"column:agent_id"`
		Count   int64  `gorm:"column:count"`
	}

	var results []countResult
	if err := r.db.WithContext(ctx).
		Model(&AgentToolPO{}).
		Select("agent_id, COUNT(*) as count").
		Where("agent_id IN ?", agentIDs).
		Group("agent_id").
		Find(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[uint64]int64, len(results))
	for _, r := range results {
		countMap[r.AgentID] = r.Count
	}
	return countMap, nil
}

func (r *AgentToolRepositoryImpl) DeleteByAgentIDs(ctx context.Context, agentIDs []uint64) error {
	if len(agentIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("agent_id IN ?", agentIDs).Delete(&AgentToolPO{}).Error
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&AgentPO{}, &AgentToolPO{})
}
