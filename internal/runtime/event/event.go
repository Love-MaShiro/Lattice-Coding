package event

import (
	"context"

	"lattice-coding/internal/common/config"
)

type EventType string

const (
	EventTypeRunCreated   EventType = "run_created"
	EventTypeRunUpdated   EventType = "run_updated"
	EventTypeRunCompleted EventType = "run_completed"
	EventTypeToolCalled   EventType = "tool_called"
	EventTypeLLMCalled    EventType = "llm_called"
)

type Event struct {
	Type      EventType
	RunID     string
	Payload   []byte
	Timestamp int64
}

type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, types []EventType, handler Handler) error
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, types []EventType, handler Handler) error
}

type Handler func(Event)

func Init(cfg *config.Config) {
}

func Publish(ctx context.Context, event Event) error {
	return nil
}

func Subscribe(ctx context.Context, types []EventType, handler Handler) error {
	return nil
}
