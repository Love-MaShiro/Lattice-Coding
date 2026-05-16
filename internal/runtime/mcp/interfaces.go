package mcp

import (
	"context"

	"lattice-coding/internal/runtime/tool"
)

type Client interface {
	ListTools(ctx context.Context) ([]ToolDefinition, error)
	CallTool(ctx context.Context, name string, input map[string]interface{}) (*ToolResult, error)
	Close() error
}

type Transport interface {
	Send(ctx context.Context, req JSONRPCRequest) (*JSONRPCResponse, error)
	Close() error
}

type ToolAdapter interface {
	Adapt(def ToolDefinition, client Client) tool.Tool
}

type ToolDefinition struct {
	Name        string
	Description string
	Schema      map[string]interface{}
}

type ToolResult struct {
	Content  string
	Data     map[string]interface{}
	IsError  bool
	Metadata map[string]interface{}
}

type JSONRPCRequest struct {
	ID     string
	Method string
	Params map[string]interface{}
}

type JSONRPCResponse struct {
	ID     string
	Result map[string]interface{}
	Error  *JSONRPCError
}

type JSONRPCError struct {
	Code    int
	Message string
	Data    map[string]interface{}
}
