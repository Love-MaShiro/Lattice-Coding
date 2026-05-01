package llm

import (
	"context"
	"errors"
	"time"

	"lattice-coding/internal/common/config"
)

var ErrPoolFull = errors.New("llm pool is full")

type SemaphorePool struct {
	sem chan struct{}
}

func NewSemaphorePool(maxConcurrent int) *SemaphorePool {
	return &SemaphorePool{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func (p *SemaphorePool) Acquire(ctx context.Context, acquireTimeout time.Duration) error {
	timer := time.NewTimer(acquireTimeout)
	defer timer.Stop()

	select {
	case p.sem <- struct{}{}:
		return nil
	case <-timer.C:
		return ErrPoolFull
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *SemaphorePool) Release() {
	select {
	case <-p.sem:
	default:
	}
}

type PoolConfig struct {
	MaxConcurrent   int
	AcquireTimeout  time.Duration
}

type LLMExecutor struct {
	Pool            *SemaphorePool
	StreamPool      *SemaphorePool
	PoolConfig      PoolConfig
	StreamPoolConfig PoolConfig
}

func NewLLMExecutor(poolCfg, streamCfg PoolConfig) *LLMExecutor {
	return &LLMExecutor{
		Pool:       NewSemaphorePool(poolCfg.MaxConcurrent),
		StreamPool: NewSemaphorePool(streamCfg.MaxConcurrent),
		PoolConfig: poolCfg,
		StreamPoolConfig: streamCfg,
	}
}

type ChatRequest struct {
	Provider string
	Model    string
	Messages []Message
}

type ChatResponse struct {
	Content string
}

type Message struct {
	Role    string
	Content string
}

type ChatClient interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

type StreamChunk struct {
	Content string
	Done    bool
	Err     error
}

type StreamClient interface {
	Stream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
}

func (e *LLMExecutor) Chat(ctx context.Context, client ChatClient, req ChatRequest) (*ChatResponse, error) {
	if err := e.Pool.Acquire(ctx, e.PoolConfig.AcquireTimeout); err != nil {
		return nil, err
	}
	defer e.Pool.Release()

	callCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return client.Chat(callCtx, req)
}

func (e *LLMExecutor) Stream(ctx context.Context, client StreamClient, req ChatRequest) (<-chan StreamChunk, error) {
	if err := e.StreamPool.Acquire(ctx, e.StreamPoolConfig.AcquireTimeout); err != nil {
		return nil, err
	}

	out := make(chan StreamChunk, 32)

	go func() {
		defer e.StreamPool.Release()
		defer close(out)

		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch, err := client.Stream(streamCtx, req)
		if err != nil {
			out <- StreamChunk{Err: err}
			return
		}

		idleTimer := time.NewTimer(120 * time.Second)
		defer idleTimer.Stop()

		for {
			select {
			case chunk, ok := <-ch:
				if !ok {
					out <- StreamChunk{Done: true}
					return
				}

				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(120 * time.Second)

				select {
				case out <- chunk:
				case <-ctx.Done():
					out <- StreamChunk{Err: ctx.Err()}
					return
				}

			case <-idleTimer.C:
				out <- StreamChunk{Err: context.DeadlineExceeded}
				return

			case <-ctx.Done():
				out <- StreamChunk{Err: ctx.Err()}
				return
			}
		}
	}()

	return out, nil
}

func Init(cfg *config.Config) {
}
