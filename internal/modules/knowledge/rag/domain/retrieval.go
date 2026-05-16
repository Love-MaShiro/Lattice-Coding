package domain

import "time"

type RetrievalRequest struct {
	Query     string         `json:"query"`
	Route     RetrievalRoute `json:"route"`
	SourceIDs []uint64       `json:"source_ids"`
	Limit     int            `json:"limit"`
	Filters   string         `json:"filters"`
}

type RetrievalResult struct {
	Query    string          `json:"query"`
	Route    RetrievalRoute  `json:"route"`
	Evidence []Evidence      `json:"evidence"`
	Trace    *RetrievalTrace `json:"trace,omitempty"`
	Metadata string          `json:"metadata"`
}

type RetrievalTrace struct {
	ID           uint64         `json:"id"`
	Query        string         `json:"query"`
	Route        RetrievalRoute `json:"route"`
	VectorCount  int            `json:"vector_count"`
	KeywordCount int            `json:"keyword_count"`
	HybridCount  int            `json:"hybrid_count"`
	LatencyMs    int64          `json:"latency_ms"`
	Metadata     string         `json:"metadata"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
