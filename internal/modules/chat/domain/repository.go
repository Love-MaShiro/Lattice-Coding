package domain

import "context"

type SessionRepository interface {
	Create(ctx context.Context, session *ChatSession) error
	Update(ctx context.Context, session *ChatSession) error
	FindByID(ctx context.Context, id uint64) (*ChatSession, error)
	FindPage(ctx context.Context, req *PageRequest) (*PageResult[*ChatSession], error)
	DeleteByID(ctx context.Context, id uint64) error
	UpdateSummary(ctx context.Context, id uint64, summary string, summarizedUntilMessageID uint64) error
}

type MessageRepository interface {
	Create(ctx context.Context, message *ChatMessage) error
	FindBySessionID(ctx context.Context, sessionID uint64, limit int) ([]*ChatMessage, error)
	FindBySessionIDAfterID(ctx context.Context, sessionID uint64, afterID uint64, limit int) ([]*ChatMessage, error)
	FindBySessionIDBeforeID(ctx context.Context, sessionID uint64, beforeID uint64, limit int) ([]*ChatMessage, error)
	CountBySessionID(ctx context.Context, sessionID uint64) (int64, error)
}
