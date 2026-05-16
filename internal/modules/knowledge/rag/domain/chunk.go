package domain

import "time"

type Chunk struct {
	ID            uint64        `json:"id"`
	ChunkKey      string        `json:"chunk_key"`
	DocumentID    uint64        `json:"document_id"`
	SourceID      uint64        `json:"source_id"`
	Ordinal       int           `json:"ordinal"`
	Title         string        `json:"title"`
	Content       string        `json:"content"`
	ContentHash   string        `json:"content_hash"`
	Location      string        `json:"location"`
	TokenEstimate int           `json:"token_estimate"`
	Strategy      ChunkStrategy `json:"strategy"`
	Metadata      string        `json:"metadata"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
