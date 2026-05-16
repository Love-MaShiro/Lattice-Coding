package domain

import "time"

type ParsedDocument struct {
	ID          uint64      `json:"id"`
	DocumentID  uint64      `json:"document_id"`
	ParserName  string      `json:"parser_name"`
	ContentText string      `json:"content_text"`
	Structure   string      `json:"structure"`
	Sections    string      `json:"sections"`
	Status      ParseStatus `json:"status"`
	Error       string      `json:"error"`
	Metadata    string      `json:"metadata"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
