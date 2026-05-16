package application

import (
	"time"

	"lattice-coding/internal/modules/workflow/domain"
)

type WorkflowDTO struct {
	ID          uint64                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      domain.WorkflowStatus `json:"status"`
	Version     int                   `json:"version"`
	Meta        string                `json:"meta"`
	Nodes       []NodeDTO             `json:"nodes,omitempty"`
	Edges       []EdgeDTO             `json:"edges,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type NodeDTO struct {
	ID         uint64            `json:"id,omitempty"`
	WorkflowID uint64            `json:"workflow_id,omitempty"`
	NodeKey    string            `json:"node_key"`
	Name       string            `json:"name"`
	Type       domain.NodeType   `json:"type"`
	Config     domain.NodeConfig `json:"config"`
	Position   string            `json:"position"`
	SortOrder  int               `json:"sort_order"`
	Meta       string            `json:"meta"`
	CreatedAt  time.Time         `json:"created_at,omitempty"`
	UpdatedAt  time.Time         `json:"updated_at,omitempty"`
}

type EdgeDTO struct {
	ID         uint64    `json:"id,omitempty"`
	WorkflowID uint64    `json:"workflow_id,omitempty"`
	EdgeKey    string    `json:"edge_key"`
	SourceKey  string    `json:"source_key"`
	TargetKey  string    `json:"target_key"`
	Condition  string    `json:"condition"`
	SortOrder  int       `json:"sort_order"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type WorkflowPageDTO struct {
	Items    []*WorkflowDTO `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type SaveWorkflowCommand struct {
	Name        string
	Description string
	Status      domain.WorkflowStatus
	Version     int
	Meta        string
	Nodes       []NodeCommand
	Edges       []EdgeCommand
}

type NodeCommand struct {
	NodeKey   string
	Name      string
	Type      domain.NodeType
	ConfigRaw string
	Position  string
	SortOrder int
	Meta      string
}

type EdgeCommand struct {
	EdgeKey   string
	SourceKey string
	TargetKey string
	Condition string
	SortOrder int
	Meta      string
}

type WorkflowPageQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Status   domain.WorkflowStatus
}
