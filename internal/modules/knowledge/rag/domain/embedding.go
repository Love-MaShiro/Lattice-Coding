package domain

import "time"

type Embedding struct {
	ID         uint64          `json:"id"`
	ChunkID    uint64          `json:"chunk_id"`
	Model      string          `json:"model"`
	Dimension  int             `json:"dimension"`
	VectorRef  string          `json:"vector_ref"`
	VectorHash string          `json:"vector_hash"`
	Status     EmbeddingStatus `json:"status"`
	Error      string          `json:"error"`
	Metadata   string          `json:"metadata"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
