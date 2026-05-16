package domain

import "time"

type Document struct {
	ID          uint64         `json:"id"`
	DocumentKey string         `json:"document_key"`
	SourceID    uint64         `json:"source_id"`
	Title       string         `json:"title"`
	Type        DocumentType   `json:"type"`
	URI         string         `json:"uri"`
	Author      string         `json:"author"`
	Version     string         `json:"version"`
	Summary     string         `json:"summary"`
	ContentHash string         `json:"content_hash"`
	Status      DocumentStatus `json:"status"`
	ParserName  string         `json:"parser_name"`
	ChunkCount  int            `json:"chunk_count"`
	Metadata    string         `json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
