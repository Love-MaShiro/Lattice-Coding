package tool

import (
	"context"
	"time"
)

type ToolType string

const (
	ToolTypeFileRead        ToolType = "file_read"
	ToolTypeFileWrite       ToolType = "file_write"
	ToolTypeFileEdit        ToolType = "file_edit"
	ToolTypeFileList        ToolType = "file_list"
	ToolTypeCodeGrep        ToolType = "code_grep"
	ToolTypeShellRun        ToolType = "shell_run"
	ToolTypeGitStatus       ToolType = "git_status"
	ToolTypeGitDiff         ToolType = "git_diff"
	ToolTypeKnowledgeSearch ToolType = "knowledge_search"
	ToolTypeEvidenceJudge   ToolType = "evidence_judge"
)

type PermissionDecision string

const (
	PermissionAllow PermissionDecision = "allow"
	PermissionDeny  PermissionDecision = "deny"
	PermissionAsk   PermissionDecision = "ask"
)

type ToolDescriptor struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Prompt          string                 `json:"prompt,omitempty"`
	Schema          Schema                 `json:"schema,omitempty"`
	ReadOnly        bool                   `json:"read_only"`
	ConcurrencySafe bool                   `json:"concurrency_safe"`
	Destructive     bool                   `json:"destructive"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type ToolContext struct {
	RunID      string                 `json:"run_id,omitempty"`
	AgentID    string                 `json:"agent_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	ProjectID  string                 `json:"project_id,omitempty"`
	SessionID  string                 `json:"session_id,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	WorkingDir string                 `json:"working_dir,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type ToolRequest struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name"`
	Input   map[string]interface{} `json:"input,omitempty"`
	Context ToolContext            `json:"context,omitempty"`
}

type ToolOutput struct {
	Content   string                 `json:"content,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Truncated bool                   `json:"truncated,omitempty"`
}

type ToolResult struct {
	RequestID     string                 `json:"request_id,omitempty"`
	ToolName      string                 `json:"tool_name"`
	IsError       bool                   `json:"is_error"`
	Content       string                 `json:"content,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Error         string                 `json:"error,omitempty"`
	FullResultRef string                 `json:"full_result_ref,omitempty"`
	StartedAt     time.Time              `json:"started_at"`
	FinishedAt    time.Time              `json:"finished_at"`
	DurationMs    int64                  `json:"duration_ms"`
	Truncated     bool                   `json:"truncated,omitempty"`
}

type Tool interface {
	Name() string
	Description() string
	Prompt() string
	Schema() Schema
	Validate(ctx context.Context, input map[string]interface{}) error
	IsReadOnly() bool
	IsConcurrencySafe() bool
	IsDestructive() bool
	CheckPermission(ctx context.Context, req ToolRequest) (PermissionDecision, string, error)
	Execute(ctx context.Context, req ToolRequest) (ToolOutput, error)
}

type BaseTool struct{}

func (BaseTool) Prompt() string { return "" }

func (BaseTool) Schema() Schema { return nil }

func (BaseTool) Validate(context.Context, map[string]interface{}) error { return nil }

func (BaseTool) IsReadOnly() bool { return false }

func (BaseTool) IsConcurrencySafe() bool { return false }

func (BaseTool) IsDestructive() bool { return false }

func (BaseTool) CheckPermission(context.Context, ToolRequest) (PermissionDecision, string, error) {
	return PermissionDeny, "tool permission check is not implemented", nil
}
