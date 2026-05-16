package domain

import "time"

type Source struct {
	ID          uint64       `json:"id"`
	SourceKey   string       `json:"source_key"`
	Name        string       `json:"name"`
	Type        SourceType   `json:"type"`
	URI         string       `json:"uri"`
	Description string       `json:"description"`
	Owner       string       `json:"owner"`
	Status      SourceStatus `json:"status"`
	Config      string       `json:"config"`
	Metadata    string       `json:"metadata"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
