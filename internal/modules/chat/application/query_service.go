package application

import (
	"context"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/chat/domain"
)

type QueryService struct {
	sessionRepo domain.SessionRepository
	messageRepo domain.MessageRepository
}

func NewQueryService(sessionRepo domain.SessionRepository, messageRepo domain.MessageRepository) *QueryService {
	return &QueryService{
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
	}
}

func (s *QueryService) GetSession(ctx context.Context, id uint64) (*SessionDTO, error) {
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询会话失败")
	}
	if session == nil {
		return nil, errors.NotFoundErr("会话不存在")
	}
	return ToSessionDTO(session), nil
}

func (s *QueryService) ListSessions(ctx context.Context, query *SessionPageQuery) (*domain.PageResult[*SessionDTO], error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.sessionRepo.FindPage(ctx, &domain.PageRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询会话列表失败")
	}

	return &domain.PageResult[*SessionDTO]{
		Items:    ToSessionDTOs(result.Items),
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *QueryService) ListMessages(ctx context.Context, query *MessageQuery) ([]*MessageDTO, error) {
	if query.SessionID == 0 {
		return nil, errors.InvalidArg("session_id is required")
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}

	messages, err := s.messageRepo.FindBySessionID(ctx, query.SessionID, limit)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询消息列表失败")
	}
	return ToMessageDTOs(messages), nil
}
