package persistence

import contextdomain "lattice-coding/internal/modules/knowledge/context/domain"

func sourceToPO(source *contextdomain.ContextSource, po *ContextSourcePO) {
	po.ID = source.ID
	po.SourceKey = source.SourceKey
	po.Kind = string(source.Kind)
	po.Name = source.Name
	po.URI = source.URI
	po.Scope = source.Scope
	po.Metadata = source.Metadata
}

func poToSource(po *ContextSourcePO) *contextdomain.ContextSource {
	return &contextdomain.ContextSource{
		ID:        po.ID,
		SourceKey: po.SourceKey,
		Kind:      contextdomain.ContextSourceKind(po.Kind),
		Name:      po.Name,
		URI:       po.URI,
		Scope:     po.Scope,
		Metadata:  po.Metadata,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func candidateToPO(candidate *contextdomain.ContextCandidate, po *ContextCandidatePO) {
	po.ID = candidate.ID
	po.CandidateKey = candidate.CandidateKey
	po.SourceKey = candidate.SourceKey
	po.SourceKind = string(candidate.SourceKind)
	po.Title = candidate.Title
	po.Content = candidate.Content
	po.Location = candidate.Location
	po.Score = candidate.Score
	po.TokenEstimate = candidate.TokenEstimate
	po.Status = string(candidate.Status)
	po.Metadata = candidate.Metadata
}

func poToCandidate(po *ContextCandidatePO) *contextdomain.ContextCandidate {
	return &contextdomain.ContextCandidate{
		ID:            po.ID,
		CandidateKey:  po.CandidateKey,
		SourceKey:     po.SourceKey,
		SourceKind:    contextdomain.ContextSourceKind(po.SourceKind),
		Title:         po.Title,
		Content:       po.Content,
		Location:      po.Location,
		Score:         po.Score,
		TokenEstimate: po.TokenEstimate,
		Status:        contextdomain.ContextCandidateStatus(po.Status),
		Metadata:      po.Metadata,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

func signalToPO(candidateID uint64, signal *contextdomain.ContextSignal, po *ContextSignalPO) {
	po.ID = signal.ID
	po.CandidateID = candidateID
	po.SignalKey = signal.SignalKey
	po.Kind = string(signal.Kind)
	po.Weight = signal.Weight
	po.Reason = signal.Reason
	po.Metadata = signal.Metadata
}

func poToSignal(po *ContextSignalPO) contextdomain.ContextSignal {
	return contextdomain.ContextSignal{
		ID:          po.ID,
		CandidateID: po.CandidateID,
		SignalKey:   po.SignalKey,
		Kind:        contextdomain.ContextSignalKind(po.Kind),
		Weight:      po.Weight,
		Reason:      po.Reason,
		Metadata:    po.Metadata,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func policyToPO(policy *contextdomain.ContextPolicy, po *ContextPolicyPO) {
	po.ID = policy.ID
	po.PolicyKey = policy.PolicyKey
	po.Name = policy.Name
	po.Description = policy.Description
	po.MaxTokens = policy.MaxTokens
	po.MaxItems = policy.MaxItems
	po.Rules = policy.Rules
	po.Metadata = policy.Metadata
}

func poToPolicy(po *ContextPolicyPO) *contextdomain.ContextPolicy {
	return &contextdomain.ContextPolicy{
		ID:          po.ID,
		PolicyKey:   po.PolicyKey,
		Name:        po.Name,
		Description: po.Description,
		MaxTokens:   po.MaxTokens,
		MaxItems:    po.MaxItems,
		Rules:       po.Rules,
		Metadata:    po.Metadata,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}
