package persistence

import (
	"context"

	"lattice-coding/internal/modules/chat/domain"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepositoryImpl(db *gorm.DB) domain.SessionRepository {
	return &SessionRepositoryImpl{db: db}
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, session *domain.ChatSession) error {
	po := &ChatSessionPO{}
	ConvertSessionToPO(session, po)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*session = *ConvertPOToSession(po)
	return nil
}

func (r *SessionRepositoryImpl) Update(ctx context.Context, session *domain.ChatSession) error {
	po := &ChatSessionPO{}
	ConvertSessionToPO(session, po)
	return r.db.WithContext(ctx).Model(po).Omit("created_at").Updates(po).Error
}

func (r *SessionRepositoryImpl) FindByID(ctx context.Context, id uint64) (*domain.ChatSession, error) {
	var po ChatSessionPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		return nil, err
	}
	return ConvertPOToSession(&po), nil
}

func (r *SessionRepositoryImpl) FindPage(ctx context.Context, req *domain.PageRequest) (*domain.PageResult[*domain.ChatSession], error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ChatSessionPO{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	var pos []ChatSessionPO
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&pos).Error; err != nil {
		return nil, err
	}

	items := make([]*domain.ChatSession, len(pos))
	for i := range pos {
		items[i] = ConvertPOToSession(&pos[i])
	}

	return &domain.PageResult[*domain.ChatSession]{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *SessionRepositoryImpl) DeleteByID(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&ChatSessionPO{}, id).Error
}

func (r *SessionRepositoryImpl) UpdateSummary(ctx context.Context, id uint64, summary string, summarizedUntilMessageID uint64) error {
	return r.db.WithContext(ctx).
		Model(&ChatSessionPO{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"summary":                     summary,
			"summarized_until_message_id": summarizedUntilMessageID,
		}).Error
}

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepositoryImpl(db *gorm.DB) domain.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, message *domain.ChatMessage) error {
	po := &ChatMessagePO{}
	ConvertMessageToPO(message, po)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	*message = *ConvertPOToMessage(po)
	return nil
}

func (r *MessageRepositoryImpl) FindBySessionID(ctx context.Context, sessionID uint64, limit int) ([]*domain.ChatMessage, error) {
	var pos []ChatMessagePO
	query := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&pos).Error; err != nil {
		return nil, err
	}
	reverseMessages(pos)
	return ConvertPOsToMessages(pos), nil
}

func (r *MessageRepositoryImpl) FindBySessionIDAfterID(ctx context.Context, sessionID uint64, afterID uint64, limit int) ([]*domain.ChatMessage, error) {
	var pos []ChatMessagePO
	query := r.db.WithContext(ctx).Where("session_id = ?", sessionID)
	if afterID > 0 {
		query = query.Where("id > ?", afterID)
	}
	query = query.Order("id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&pos).Error; err != nil {
		return nil, err
	}
	return ConvertPOsToMessages(pos), nil
}
func (r *MessageRepositoryImpl) FindBySessionIDBeforeID(ctx context.Context, sessionID uint64, beforeID uint64, limit int) ([]*domain.ChatMessage, error) {
	var pos []ChatMessagePO
	query := r.db.WithContext(ctx).Where("session_id = ?", sessionID)
	if beforeID > 0 {
		query = query.Where("id < ?", beforeID)
	}
	query = query.Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&pos).Error; err != nil {
		return nil, err
	}
	reverseMessages(pos)
	return ConvertPOsToMessages(pos), nil
}

func (r *MessageRepositoryImpl) CountBySessionID(ctx context.Context, sessionID uint64) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&ChatMessagePO{}).Where("session_id = ?", sessionID).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&ChatSessionPO{}, &ChatMessagePO{})
}

func reverseMessages(pos []ChatMessagePO) {
	for i, j := 0, len(pos)-1; i < j; i, j = i+1, j-1 {
		pos[i], pos[j] = pos[j], pos[i]
	}
}
