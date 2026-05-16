package domain

import "time"

type KeywordIndex struct {
	ID         uint64             `json:"id"`
	ChunkID    uint64             `json:"chunk_id"`
	IndexName  string             `json:"index_name"`
	DocumentID uint64             `json:"document_id"`
	SourceID   uint64             `json:"source_id"`
	Status     KeywordIndexStatus `json:"status"`
	IndexedAt  *time.Time         `json:"indexed_at,omitempty"`
	Error      string             `json:"error"`
	Metadata   string             `json:"metadata"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
