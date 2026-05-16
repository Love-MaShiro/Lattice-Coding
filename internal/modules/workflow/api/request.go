package api

import (
	"encoding/json"

	"lattice-coding/internal/modules/workflow/application"
	"lattice-coding/internal/modules/workflow/domain"
)

type WorkflowPageQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   string `form:"status" json:"status"`
}

type SaveWorkflowRequest struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Version     int           `json:"version"`
	Meta        string        `json:"meta"`
	Nodes       []NodeRequest `json:"nodes"`
	Edges       []EdgeRequest `json:"edges"`
}

type NodeRequest struct {
	NodeKey   string          `json:"node_key" binding:"required"`
	Name      string          `json:"name"`
	Type      string          `json:"type" binding:"required"`
	Config    json.RawMessage `json:"config"`
	Position  string          `json:"position"`
	SortOrder int             `json:"sort_order"`
	Meta      string          `json:"meta"`
}

type EdgeRequest struct {
	EdgeKey   string `json:"edge_key"`
	SourceKey string `json:"source_key" binding:"required"`
	TargetKey string `json:"target_key" binding:"required"`
	Condition string `json:"condition"`
	SortOrder int    `json:"sort_order"`
	Meta      string `json:"meta"`
}

type CodeReviewRequest struct {
	RunID         string `json:"run_id"`
	AgentID       string `json:"agent_id"`
	SessionID     string `json:"session_id"`
	WorkingDir    string `json:"working_dir" binding:"required"`
	ModelConfigID uint64 `json:"model_config_id" binding:"required"`
	MaxFiles      int    `json:"max_files"`
	MaxChars      int    `json:"max_chars"`
}

func (r *WorkflowPageQuery) ToApplication() application.WorkflowPageQuery {
	return application.WorkflowPageQuery{
		Page:     r.Page,
		PageSize: r.PageSize,
		Keyword:  r.Keyword,
		Status:   domain.WorkflowStatus(r.Status),
	}
}

func (r *SaveWorkflowRequest) ToApplication() application.SaveWorkflowCommand {
	nodes := make([]application.NodeCommand, len(r.Nodes))
	for i := range r.Nodes {
		nodes[i] = application.NodeCommand{
			NodeKey:   r.Nodes[i].NodeKey,
			Name:      r.Nodes[i].Name,
			Type:      domain.NodeType(r.Nodes[i].Type),
			ConfigRaw: rawConfigString(r.Nodes[i].Config),
			Position:  r.Nodes[i].Position,
			SortOrder: r.Nodes[i].SortOrder,
			Meta:      r.Nodes[i].Meta,
		}
	}
	edges := make([]application.EdgeCommand, len(r.Edges))
	for i := range r.Edges {
		edges[i] = application.EdgeCommand{
			EdgeKey:   r.Edges[i].EdgeKey,
			SourceKey: r.Edges[i].SourceKey,
			TargetKey: r.Edges[i].TargetKey,
			Condition: r.Edges[i].Condition,
			SortOrder: r.Edges[i].SortOrder,
			Meta:      r.Edges[i].Meta,
		}
	}
	return application.SaveWorkflowCommand{
		Name:        r.Name,
		Description: r.Description,
		Status:      domain.WorkflowStatus(r.Status),
		Version:     r.Version,
		Meta:        r.Meta,
		Nodes:       nodes,
		Edges:       edges,
	}
}

func (r *CodeReviewRequest) ToApplication() application.CodeReviewCommand {
	return application.CodeReviewCommand{
		RunID:         r.RunID,
		AgentID:       r.AgentID,
		SessionID:     r.SessionID,
		WorkingDir:    r.WorkingDir,
		ModelConfigID: r.ModelConfigID,
		MaxFiles:      r.MaxFiles,
		MaxChars:      r.MaxChars,
	}
}

func rawConfigString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "{}"
	}
	return string(raw)
}
