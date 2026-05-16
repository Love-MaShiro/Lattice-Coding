package domain

import "time"

type KnowledgePack struct {
	ID            uint64              `json:"id"`
	PackKey       string              `json:"pack_key"`
	Query         string              `json:"query"`
	Intent        string              `json:"intent"`
	Route         string              `json:"route"`
	Status        PackStatus          `json:"status"`
	TokenEstimate int                 `json:"token_estimate"`
	MaxTokens     int                 `json:"max_tokens"`
	PromptContext string              `json:"prompt_context"`
	Warnings      string              `json:"warnings"`
	Options       string              `json:"options"`
	Meta          string              `json:"meta"`
	Items         []KnowledgeItem     `json:"items"`
	Citations     []KnowledgeCitation `json:"citations"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type KnowledgeItem struct {
	ID            uint64         `json:"id"`
	PackID        uint64         `json:"pack_id"`
	ItemKey       string         `json:"item_key"`
	SourceKind    PackSourceKind `json:"source_kind"`
	SourceID      string         `json:"source_id"`
	SourceType    string         `json:"source_type"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Location      string         `json:"location"`
	Score         float64        `json:"score"`
	TokenEstimate int            `json:"token_estimate"`
	CitationKey   string         `json:"citation_key"`
	Metadata      string         `json:"metadata"`
	SortOrder     int            `json:"sort_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type KnowledgeCitation struct {
	ID          uint64         `json:"id"`
	PackID      uint64         `json:"pack_id"`
	CitationKey string         `json:"citation_key"`
	SourceKind  PackSourceKind `json:"source_kind"`
	SourceID    string         `json:"source_id"`
	Title       string         `json:"title"`
	Location    string         `json:"location"`
	URI         string         `json:"uri"`
	Score       float64        `json:"score"`
	Metadata    string         `json:"metadata"`
	SortOrder   int            `json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type PackCandidate struct {
	CandidateKey  string         `json:"candidate_key"`
	SourceKind    PackSourceKind `json:"source_kind"`
	SourceID      string         `json:"source_id"`
	SourceType    string         `json:"source_type"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Location      string         `json:"location"`
	URI           string         `json:"uri"`
	Score         float64        `json:"score"`
	TokenEstimate int            `json:"token_estimate"`
	Metadata      string         `json:"metadata"`
}
