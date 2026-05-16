package domain

import "time"

type EvidenceGateRequest struct {
	Query       string     `json:"query"`
	AnswerDraft string     `json:"answer_draft"`
	Evidence    []Evidence `json:"evidence"`
	Policy      string     `json:"policy"`
}

type EvidenceGateResult struct {
	ID              uint64             `json:"id"`
	Query           string             `json:"query"`
	Passed          bool               `json:"passed"`
	Score           float64            `json:"score"`
	Reason          string             `json:"reason"`
	MissingAspects  string             `json:"missing_aspects"`
	WeakEvidenceIDs string             `json:"weak_evidence_ids"`
	Contradictions  string             `json:"contradictions"`
	Suggestions     string             `json:"suggestions"`
	Status          EvidenceGateStatus `json:"status"`
	Metadata        string             `json:"metadata"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
