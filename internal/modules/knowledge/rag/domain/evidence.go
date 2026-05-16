package domain

import "time"

type Evidence struct {
	ID            uint64           `json:"id"`
	EvidenceKey   string           `json:"evidence_key"`
	Query         string           `json:"query"`
	ChunkID       uint64           `json:"chunk_id"`
	DocumentID    uint64           `json:"document_id"`
	SourceID      uint64           `json:"source_id"`
	Channel       RetrievalChannel `json:"channel"`
	Title         string           `json:"title"`
	Content       string           `json:"content"`
	Location      string           `json:"location"`
	Score         float64          `json:"score"`
	RerankScore   float64          `json:"rerank_score"`
	TokenEstimate int              `json:"token_estimate"`
	CitationURI   string           `json:"citation_uri"`
	Metadata      string           `json:"metadata"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
