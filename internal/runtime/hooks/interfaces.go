package hooks

import "context"

type EventType string

const (
	EventBeforeRun  EventType = "before_run"
	EventAfterRun   EventType = "after_run"
	EventBeforeTool EventType = "before_tool"
	EventAfterTool  EventType = "after_tool"
	EventOnError    EventType = "on_error"
)

type Engine interface {
	Execute(ctx context.Context, event Event) error
}

type Handler interface {
	Handle(ctx context.Context, event Event) error
}

type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

type Event struct {
	Type     EventType
	RunID    string
	Payload  map[string]interface{}
	Metadata map[string]interface{}
}
