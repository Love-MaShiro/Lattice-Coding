package tool

import (
	"context"
	"lattice-coding/internal/common/config"
)

type ToolType string

const (
	ToolTypeFileRead     ToolType = "file_read"
	ToolTypeFileWrite    ToolType = "file_write"
	ToolTypeCommand      ToolType = "command"
	ToolTypeGitDiff      ToolType = "git_diff"
	ToolTypeTest         ToolType = "test"
	ToolTypeSearch       ToolType = "search"
	ToolTypeKnowledge    ToolType = "knowledge"
)

type ToolRequest struct {
	Type ToolType
	Args map[string]interface{}
}

type ToolResponse struct {
	Content string
	Error   error
}

type ToolExecutor interface {
	Execute(ctx context.Context, req ToolRequest) (*ToolResponse, error)
}

func Init(cfg *config.Config) {
}

func Execute(ctx context.Context, req ToolRequest) (*ToolResponse, error) {
	return nil, nil
}
