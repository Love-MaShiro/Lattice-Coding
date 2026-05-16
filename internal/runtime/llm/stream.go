package llm

import "context"

type StreamReader <-chan StreamChunk

type StreamHandler interface {
	HandleChunk(ctx context.Context, chunk StreamChunk) error
}

type StreamHandlerFunc func(ctx context.Context, chunk StreamChunk) error

func (f StreamHandlerFunc) HandleChunk(ctx context.Context, chunk StreamChunk) error {
	return f(ctx, chunk)
}
