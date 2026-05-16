package application

import (
	"context"
	"strings"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/knowledge/context/domain"
)

type SourceService interface {
	CreateSource(ctx context.Context, source *domain.ContextSource) (*domain.ContextSource, error)
	GetSource(ctx context.Context, sourceKey string) (*domain.ContextSource, error)
}

type sourceService struct {
	repo domain.SourceRepository
}

func NewSourceService(repo domain.SourceRepository) SourceService {
	return &sourceService{repo: repo}
}

func (s *sourceService) CreateSource(ctx context.Context, source *domain.ContextSource) (*domain.ContextSource, error) {
	if source == nil {
		return nil, errors.InvalidArg("source is required")
	}
	source.SourceKey = strings.TrimSpace(source.SourceKey)
	if source.SourceKey == "" {
		return nil, errors.InvalidArg("source_key is required")
	}
	if source.Kind == "" {
		return nil, errors.InvalidArg("source kind is required")
	}
	if source.Metadata == "" {
		source.Metadata = "{}"
	}
	if err := s.repo.Create(ctx, source); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create context source failed")
	}
	return source, nil
}

func (s *sourceService) GetSource(ctx context.Context, sourceKey string) (*domain.ContextSource, error) {
	sourceKey = strings.TrimSpace(sourceKey)
	if sourceKey == "" {
		return nil, errors.InvalidArg("source_key is required")
	}
	source, err := s.repo.FindByKey(ctx, sourceKey)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "get context source failed")
	}
	return source, nil
}
