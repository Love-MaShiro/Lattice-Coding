package application

import (
	"context"
	"strings"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/knowledge/context/domain"
)

type PolicyService interface {
	SavePolicy(ctx context.Context, policy *domain.ContextPolicy) (*domain.ContextPolicy, error)
	GetPolicy(ctx context.Context, policyKey string) (*domain.ContextPolicy, error)
}

type policyService struct {
	repo domain.PolicyRepository
}

func NewPolicyService(repo domain.PolicyRepository) PolicyService {
	return &policyService{repo: repo}
}

func (s *policyService) SavePolicy(ctx context.Context, policy *domain.ContextPolicy) (*domain.ContextPolicy, error) {
	if policy == nil {
		return nil, errors.InvalidArg("policy is required")
	}
	policy.PolicyKey = strings.TrimSpace(policy.PolicyKey)
	if policy.PolicyKey == "" {
		return nil, errors.InvalidArg("policy_key is required")
	}
	if policy.MaxTokens < 0 {
		return nil, errors.InvalidArg("max_tokens must be greater than or equal to 0")
	}
	if policy.MaxItems < 0 {
		return nil, errors.InvalidArg("max_items must be greater than or equal to 0")
	}
	if policy.Rules == "" {
		policy.Rules = "{}"
	}
	if policy.Metadata == "" {
		policy.Metadata = "{}"
	}
	if policy.ID > 0 {
		if err := s.repo.Update(ctx, policy); err != nil {
			return nil, errors.DatabaseErrWithErr(err, "update context policy failed")
		}
		return policy, nil
	}
	if err := s.repo.Create(ctx, policy); err != nil {
		return nil, errors.DatabaseErrWithErr(err, "create context policy failed")
	}
	return policy, nil
}

func (s *policyService) GetPolicy(ctx context.Context, policyKey string) (*domain.ContextPolicy, error) {
	policyKey = strings.TrimSpace(policyKey)
	if policyKey == "" {
		return nil, errors.InvalidArg("policy_key is required")
	}
	policy, err := s.repo.FindByKey(ctx, policyKey)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "get context policy failed")
	}
	return policy, nil
}
