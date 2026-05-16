package domain

import "time"

type WorkflowDefinition struct {
	ID          uint64           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      WorkflowStatus   `json:"status"`
	Version     int              `json:"version"`
	Meta        string           `json:"meta"`
	Nodes       []NodeDefinition `json:"nodes"`
	Edges       []EdgeDefinition `json:"edges"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type NodeDefinition struct {
	ID         uint64     `json:"id"`
	WorkflowID uint64     `json:"workflow_id"`
	NodeKey    string     `json:"node_key"`
	Name       string     `json:"name"`
	Type       NodeType   `json:"type"`
	Config     NodeConfig `json:"config"`
	Position   string     `json:"position"`
	SortOrder  int        `json:"sort_order"`
	Meta       string     `json:"meta"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type EdgeDefinition struct {
	ID         uint64    `json:"id"`
	WorkflowID uint64    `json:"workflow_id"`
	EdgeKey    string    `json:"edge_key"`
	SourceKey  string    `json:"source_key"`
	TargetKey  string    `json:"target_key"`
	Condition  string    `json:"condition"`
	SortOrder  int       `json:"sort_order"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RunState struct {
	ID              uint64                 `json:"id"`
	WorkflowID      uint64                 `json:"workflow_id"`
	WorkflowVersion int                    `json:"workflow_version"`
	Status          RunStatus              `json:"status"`
	CurrentNodeKey  string                 `json:"current_node_key"`
	Inputs          map[string]interface{} `json:"inputs"`
	Outputs         map[string]interface{} `json:"outputs"`
	Error           string                 `json:"error"`
	StartedAt       *time.Time             `json:"started_at,omitempty"`
	FinishedAt      *time.Time             `json:"finished_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type PageRequest struct {
	Page     int
	PageSize int
	Keyword  string
	Status   WorkflowStatus
}

type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}
