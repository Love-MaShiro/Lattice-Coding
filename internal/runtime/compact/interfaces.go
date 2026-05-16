package compact

import "context"

type Policy interface {
	ShouldCompact(state State) bool
}

type Compactor interface {
	Compact(ctx context.Context, state State) (*Summary, error)
}

type Restorer interface {
	Restore(ctx context.Context, summary Summary) (*State, error)
}

type State struct {
	RunID       string
	Messages    []Message
	ToolTraces  []ToolTrace
	TokenCount  int
	TokenBudget int
	Metadata    map[string]interface{}
}

type Summary struct {
	RunID    string
	Content  string
	Metadata map[string]interface{}
}

type Message struct {
	Role    string
	Content string
}

type ToolTrace struct {
	ToolName string
	Input    map[string]interface{}
	Output   string
	IsError  bool
}
