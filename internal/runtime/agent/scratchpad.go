package agent

import (
	"context"
	"sync"
)

type Scratchpad interface {
	Append(ctx context.Context, step ReActStep) error
	Steps(ctx context.Context) ([]ReActStep, error)
}

type MemoryScratchpad struct {
	mu    sync.RWMutex
	steps []ReActStep
}

func NewMemoryScratchpad() *MemoryScratchpad {
	return &MemoryScratchpad{}
}

func (s *MemoryScratchpad) Append(ctx context.Context, step ReActStep) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.steps = append(s.steps, step)
	return nil
}

func (s *MemoryScratchpad) Steps(ctx context.Context) ([]ReActStep, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	steps := make([]ReActStep, len(s.steps))
	copy(steps, s.steps)
	return steps, nil
}
