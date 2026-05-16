package domain

import "time"

type ContextSource struct {
	ID        uint64            `json:"id"`
	SourceKey string            `json:"source_key"`
	Kind      ContextSourceKind `json:"kind"`
	Name      string            `json:"name"`
	URI       string            `json:"uri"`
	Scope     string            `json:"scope"`
	Metadata  string            `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type ContextCandidate struct {
	ID            uint64                 `json:"id"`
	CandidateKey  string                 `json:"candidate_key"`
	SourceKey     string                 `json:"source_key"`
	SourceKind    ContextSourceKind      `json:"source_kind"`
	Title         string                 `json:"title"`
	Content       string                 `json:"content"`
	Location      string                 `json:"location"`
	Score         float64                `json:"score"`
	TokenEstimate int                    `json:"token_estimate"`
	Status        ContextCandidateStatus `json:"status"`
	Signals       []ContextSignal        `json:"signals"`
	Metadata      string                 `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type ContextSignal struct {
	ID          uint64            `json:"id"`
	CandidateID uint64            `json:"candidate_id"`
	SignalKey   string            `json:"signal_key"`
	Kind        ContextSignalKind `json:"kind"`
	Weight      float64           `json:"weight"`
	Reason      string            `json:"reason"`
	Metadata    string            `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ContextPolicy struct {
	ID          uint64    `json:"id"`
	PolicyKey   string    `json:"policy_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxTokens   int       `json:"max_tokens"`
	MaxItems    int       `json:"max_items"`
	Rules       string    `json:"rules"`
	Metadata    string    `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ContextPack struct {
	PackKey       string            `json:"pack_key"`
	Query         string            `json:"query"`
	PolicyKey     string            `json:"policy_key"`
	MaxTokens     int               `json:"max_tokens"`
	TokenEstimate int               `json:"token_estimate"`
	PromptContext string            `json:"prompt_context"`
	Items         []ContextPackItem `json:"items"`
	Warnings      []string          `json:"warnings"`
	Metadata      string            `json:"metadata"`
}

type ContextPackItem struct {
	CandidateKey  string            `json:"candidate_key"`
	SourceKey     string            `json:"source_key"`
	SourceKind    ContextSourceKind `json:"source_kind"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	Location      string            `json:"location"`
	Score         float64           `json:"score"`
	TokenEstimate int               `json:"token_estimate"`
	Signals       []ContextSignal   `json:"signals"`
	Metadata      string            `json:"metadata"`
	SortOrder     int               `json:"sort_order"`
}
