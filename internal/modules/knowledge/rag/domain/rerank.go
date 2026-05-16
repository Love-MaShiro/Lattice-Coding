package domain

import "time"

type RerankRequest struct {
	Query    string     `json:"query"`
	Evidence []Evidence `json:"evidence"`
	TopK     int        `json:"top_k"`
	Policy   string     `json:"policy"`
}

type RerankResult struct {
	ID          uint64     `json:"id"`
	Query       string     `json:"query"`
	Model       string     `json:"model"`
	InputCount  int        `json:"input_count"`
	OutputCount int        `json:"output_count"`
	Items       []Evidence `json:"items"`
	Metadata    string     `json:"metadata"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
