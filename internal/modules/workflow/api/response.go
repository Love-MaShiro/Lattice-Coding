package api

import (
	"encoding/json"
	"time"

	"lattice-coding/internal/modules/workflow/application"
	"lattice-coding/internal/modules/workflow/domain"
)

type WorkflowResponse struct {
	ID          uint64         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Version     int            `json:"version"`
	Meta        string         `json:"meta"`
	Nodes       []NodeResponse `json:"nodes,omitempty"`
	Edges       []EdgeResponse `json:"edges,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type NodeResponse struct {
	ID         uint64          `json:"id,omitempty"`
	WorkflowID uint64          `json:"workflow_id,omitempty"`
	NodeKey    string          `json:"node_key"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Config     json.RawMessage `json:"config"`
	Position   string          `json:"position"`
	SortOrder  int             `json:"sort_order"`
	Meta       string          `json:"meta"`
	CreatedAt  time.Time       `json:"created_at,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at,omitempty"`
}

type EdgeResponse struct {
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

func ToWorkflowResponse(dto *application.WorkflowDTO) WorkflowResponse {
	if dto == nil {
		return WorkflowResponse{}
	}
	return WorkflowResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      string(dto.Status),
		Version:     dto.Version,
		Meta:        dto.Meta,
		Nodes:       ToNodeResponses(dto.Nodes),
		Edges:       ToEdgeResponses(dto.Edges),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func ToWorkflowResponses(dtos []*application.WorkflowDTO) []WorkflowResponse {
	items := make([]WorkflowResponse, len(dtos))
	for i := range dtos {
		items[i] = ToWorkflowResponse(dtos[i])
	}
	return items
}

func ToNodeResponses(dtos []application.NodeDTO) []NodeResponse {
	items := make([]NodeResponse, len(dtos))
	for i := range dtos {
		items[i] = NodeResponse{
			ID:         dtos[i].ID,
			WorkflowID: dtos[i].WorkflowID,
			NodeKey:    dtos[i].NodeKey,
			Name:       dtos[i].Name,
			Type:       string(dtos[i].Type),
			Config:     nodeConfigRaw(dtos[i].Config),
			Position:   dtos[i].Position,
			SortOrder:  dtos[i].SortOrder,
			Meta:       dtos[i].Meta,
			CreatedAt:  dtos[i].CreatedAt,
			UpdatedAt:  dtos[i].UpdatedAt,
		}
	}
	return items
}

func ToEdgeResponses(dtos []application.EdgeDTO) []EdgeResponse {
	items := make([]EdgeResponse, len(dtos))
	for i := range dtos {
		items[i] = EdgeResponse{
			ID:         dtos[i].ID,
			WorkflowID: dtos[i].WorkflowID,
			EdgeKey:    dtos[i].EdgeKey,
			SourceKey:  dtos[i].SourceKey,
			TargetKey:  dtos[i].TargetKey,
			Condition:  dtos[i].Condition,
			SortOrder:  dtos[i].SortOrder,
			Meta:       dtos[i].Meta,
			CreatedAt:  dtos[i].CreatedAt,
			UpdatedAt:  dtos[i].UpdatedAt,
		}
	}
	return items
}

func nodeConfigRaw(config domain.NodeConfig) json.RawMessage {
	if config == nil {
		return json.RawMessage("{}")
	}
	raw := config.RawJSON()
	if raw == "" {
		raw = "{}"
	}
	return json.RawMessage(raw)
}
