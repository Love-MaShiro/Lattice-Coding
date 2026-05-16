package application

import (
	"context"
	"sort"
	"strings"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/knowledge/context/domain"
)

const defaultContextMaxTokens = 8000

type BuildContextPackCommand struct {
	PackKey       string
	Query         string
	PolicyKey     string
	SourceKeys    []string
	CandidateKeys []string
	MaxTokens     int
	MaxItems      int
	Metadata      string
}

type ContextPackService interface {
	BuildContextPack(ctx context.Context, cmd BuildContextPackCommand) (*domain.ContextPack, error)
}

type contextPackService struct {
	candidateRepo domain.CandidateRepository
	policyRepo    domain.PolicyRepository
}

func NewContextPackService(candidateRepo domain.CandidateRepository, policyRepo domain.PolicyRepository) ContextPackService {
	return &contextPackService{
		candidateRepo: candidateRepo,
		policyRepo:    policyRepo,
	}
}

func (s *contextPackService) BuildContextPack(ctx context.Context, cmd BuildContextPackCommand) (*domain.ContextPack, error) {
	cmd.PackKey = strings.TrimSpace(cmd.PackKey)
	if cmd.PackKey == "" {
		return nil, errors.InvalidArg("pack_key is required")
	}

	policy, err := s.loadPolicy(ctx, cmd.PolicyKey)
	if err != nil {
		return nil, err
	}
	maxTokens, maxItems := resolveLimits(cmd, policy)

	candidates, err := s.loadCandidates(ctx, cmd)
	if err != nil {
		return nil, err
	}
	ranked := rankCandidates(candidates)
	items, warnings := selectContextItems(ranked, maxTokens, maxItems)

	pack := &domain.ContextPack{
		PackKey:       cmd.PackKey,
		Query:         cmd.Query,
		PolicyKey:     strings.TrimSpace(cmd.PolicyKey),
		MaxTokens:     maxTokens,
		TokenEstimate: sumItemTokens(items),
		Items:         items,
		Warnings:      warnings,
		Metadata:      firstNonEmpty(cmd.Metadata, "{}"),
	}
	pack.PromptContext = renderPromptContext(pack)
	return pack, nil
}

func (s *contextPackService) loadPolicy(ctx context.Context, policyKey string) (*domain.ContextPolicy, error) {
	policyKey = strings.TrimSpace(policyKey)
	if policyKey == "" {
		return nil, nil
	}
	policy, err := s.policyRepo.FindByKey(ctx, policyKey)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "get context policy failed")
	}
	return policy, nil
}

func (s *contextPackService) loadCandidates(ctx context.Context, cmd BuildContextPackCommand) ([]*domain.ContextCandidate, error) {
	seen := map[string]struct{}{}
	candidates := make([]*domain.ContextCandidate, 0)

	for _, candidateKey := range cmd.CandidateKeys {
		candidateKey = strings.TrimSpace(candidateKey)
		if candidateKey == "" {
			continue
		}
		if _, ok := seen[candidateKey]; ok {
			continue
		}
		candidate, err := s.candidateRepo.FindByKeyWithSignals(ctx, candidateKey)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "get context candidate failed")
		}
		candidates = append(candidates, candidate)
		seen[candidateKey] = struct{}{}
	}

	for _, sourceKey := range cmd.SourceKeys {
		sourceKey = strings.TrimSpace(sourceKey)
		if sourceKey == "" {
			continue
		}
		sourceCandidates, err := s.candidateRepo.FindBySourceKey(ctx, sourceKey)
		if err != nil {
			return nil, errors.DatabaseErrWithErr(err, "list context candidates failed")
		}
		for _, candidate := range sourceCandidates {
			if candidate == nil {
				continue
			}
			if _, ok := seen[candidate.CandidateKey]; ok {
				continue
			}
			withSignals, err := s.candidateRepo.FindByKeyWithSignals(ctx, candidate.CandidateKey)
			if err != nil {
				return nil, errors.DatabaseErrWithErr(err, "get context candidate signals failed")
			}
			candidates = append(candidates, withSignals)
			seen[candidate.CandidateKey] = struct{}{}
		}
	}

	return candidates, nil
}

func resolveLimits(cmd BuildContextPackCommand, policy *domain.ContextPolicy) (int, int) {
	maxTokens := cmd.MaxTokens
	maxItems := cmd.MaxItems
	if policy != nil {
		if maxTokens <= 0 {
			maxTokens = policy.MaxTokens
		}
		if maxItems <= 0 {
			maxItems = policy.MaxItems
		}
	}
	if maxTokens <= 0 {
		maxTokens = defaultContextMaxTokens
	}
	return maxTokens, maxItems
}

func rankCandidates(candidates []*domain.ContextCandidate) []*domain.ContextCandidate {
	ranked := make([]*domain.ContextCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate == nil || candidate.Status == domain.ContextCandidateRejected {
			continue
		}
		ranked = append(ranked, candidate)
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		return candidateRankScore(ranked[i]) > candidateRankScore(ranked[j])
	})
	return ranked
}

func candidateRankScore(candidate *domain.ContextCandidate) float64 {
	score := candidate.Score
	for _, signal := range candidate.Signals {
		score += signal.Weight
	}
	return score
}

func selectContextItems(candidates []*domain.ContextCandidate, maxTokens int, maxItems int) ([]domain.ContextPackItem, []string) {
	items := make([]domain.ContextPackItem, 0, len(candidates))
	warnings := make([]string, 0)
	usedTokens := 0

	for _, candidate := range candidates {
		if maxItems > 0 && len(items) >= maxItems {
			warnings = append(warnings, "context item limit reached")
			break
		}
		tokenEstimate := candidate.TokenEstimate
		if tokenEstimate <= 0 {
			tokenEstimate = estimateTokens(candidate.Content)
		}
		if usedTokens+tokenEstimate > maxTokens {
			warnings = append(warnings, "context token limit reached")
			continue
		}
		usedTokens += tokenEstimate
		items = append(items, domain.ContextPackItem{
			CandidateKey:  candidate.CandidateKey,
			SourceKey:     candidate.SourceKey,
			SourceKind:    candidate.SourceKind,
			Title:         candidate.Title,
			Content:       candidate.Content,
			Location:      candidate.Location,
			Score:         candidateRankScore(candidate),
			TokenEstimate: tokenEstimate,
			Signals:       candidate.Signals,
			Metadata:      candidate.Metadata,
			SortOrder:     len(items) + 1,
		})
	}

	return items, dedupeStrings(warnings)
}

func renderPromptContext(pack *domain.ContextPack) string {
	var builder strings.Builder
	builder.WriteString("# Context Pack\n\n")
	if strings.TrimSpace(pack.Query) != "" {
		builder.WriteString("Query: ")
		builder.WriteString(strings.TrimSpace(pack.Query))
		builder.WriteString("\n\n")
	}
	for _, item := range pack.Items {
		builder.WriteString("## ")
		builder.WriteString(item.Title)
		builder.WriteString("\n")
		if item.Location != "" {
			builder.WriteString("Location: ")
			builder.WriteString(item.Location)
			builder.WriteString("\n")
		}
		builder.WriteString("Source: ")
		builder.WriteString(string(item.SourceKind))
		builder.WriteString("/")
		builder.WriteString(item.SourceKey)
		builder.WriteString("\n\n")
		builder.WriteString(strings.TrimSpace(item.Content))
		builder.WriteString("\n\n")
	}
	return strings.TrimSpace(builder.String())
}

func sumItemTokens(items []domain.ContextPackItem) int {
	total := 0
	for _, item := range items {
		total += item.TokenEstimate
	}
	return total
}

func estimateTokens(content string) int {
	content = strings.TrimSpace(content)
	if content == "" {
		return 0
	}
	estimate := len([]rune(content)) / 4
	if estimate <= 0 {
		return 1
	}
	return estimate
}

func firstNonEmpty(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
