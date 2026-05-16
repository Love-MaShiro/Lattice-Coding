package application

import (
	"context"
	"strings"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/knowledge/context/domain"
)

type CandidateService interface {
	CreateCandidate(ctx context.Context, candidate *domain.ContextCandidate) (*domain.ContextCandidate, error)
	ListCandidatesBySource(ctx context.Context, sourceKey string) ([]*domain.ContextCandidate, error)
}

type candidateService struct {
	repo domain.CandidateRepository
}

func NewCandidateService(repo domain.CandidateRepository) CandidateService {
	return &candidateService{repo: repo}
}

func (s *candidateService) CreateCandidate(ctx context.Context, candidate *domain.ContextCandidate) (*domain.ContextCandidate, error) {
	if candidate == nil {
		return nil, errors.InvalidArg("candidate is required")
	}
	candidate.CandidateKey = strings.TrimSpace(candidate.CandidateKey)
	candidate.SourceKey = strings.TrimSpace(candidate.SourceKey)
	if candidate.CandidateKey == "" {
		return nil, errors.InvalidArg("candidate_key is required")
	}
	if candidate.SourceKey == "" {
		return nil, errors.InvalidArg("source_key is required")
	}
	if candidate.SourceKind == "" {
		return nil, errors.InvalidArg("source_kind is required")
	}
	if candidate.Status == "" {
		candidate.Status = domain.ContextCandidatePending
	}
	if candidate.Metadata == "" {
		candidate.Metadata = "{}"
	}
	if err := s.repo.CreateWithSignals(ctx, candidate); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create context candidate failed")
	}
	return candidate, nil
}

func (s *candidateService) ListCandidatesBySource(ctx context.Context, sourceKey string) ([]*domain.ContextCandidate, error) {
	sourceKey = strings.TrimSpace(sourceKey)
	if sourceKey == "" {
		return nil, errors.InvalidArg("source_key is required")
	}
	candidates, err := s.repo.FindBySourceKey(ctx, sourceKey)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "list context candidates failed")
	}
	return candidates, nil
}
