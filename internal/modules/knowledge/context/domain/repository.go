package domain

import "context"

type SourceRepository interface {
	Create(ctx context.Context, source *ContextSource) error
	FindByKey(ctx context.Context, sourceKey string) (*ContextSource, error)
	DeleteByKey(ctx context.Context, sourceKey string) error
}

type CandidateRepository interface {
	CreateWithSignals(ctx context.Context, candidate *ContextCandidate) error
	FindByKeyWithSignals(ctx context.Context, candidateKey string) (*ContextCandidate, error)
	FindBySourceKey(ctx context.Context, sourceKey string) ([]*ContextCandidate, error)
	DeleteByKey(ctx context.Context, candidateKey string) error
}

type PolicyRepository interface {
	Create(ctx context.Context, policy *ContextPolicy) error
	Update(ctx context.Context, policy *ContextPolicy) error
	FindByKey(ctx context.Context, policyKey string) (*ContextPolicy, error)
	DeleteByKey(ctx context.Context, policyKey string) error
}
