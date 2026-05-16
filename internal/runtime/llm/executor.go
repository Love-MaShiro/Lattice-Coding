package llm

import (
	"context"
	"sync"
	"time"

	"lattice-coding/internal/common/config"
)

type Executor struct {
	pool           *SemaphorePool
	streamPool     *SemaphorePool
	retryCfg       RetryConfig
	timeoutCfg     TimeoutConfig
	breaker        *CircuitBreaker
	router         *Router
	clientRegistry *ClientRegistry
	modelFactory   *LLMFactory
	mu             sync.RWMutex
}

type TimeoutConfig struct {
	SyncCall    time.Duration
	StreamIdle  time.Duration
	HealthCheck time.Duration
}

func NewExecutor(cfg *config.LLMConfig) *Executor {
	poolCfg := cfg.Pool
	streamCfg := cfg.Stream

	executor := &Executor{
		pool:       NewSemaphorePool(poolCfg.MaxConcurrent),
		streamPool: NewSemaphorePool(streamCfg.MaxConcurrent),
		retryCfg: RetryConfig{
			MaxAttempts:     2,
			InitialInterval: 500 * time.Millisecond,
			MaxInterval:     5 * time.Second,
		},
		timeoutCfg: TimeoutConfig{
			SyncCall:    60 * time.Second,
			StreamIdle:  120 * time.Second,
			HealthCheck: 10 * time.Second,
		},
		clientRegistry: NewClientRegistry(),
	}
	executor.router = NewRouter(executor.clientRegistry, RouterConfig{})

	if cfg.CircuitBreaker.Enabled {
		breakerCfg := BreakerConfig{
			Name:                 "llm",
			FailureRateThreshold: cfg.CircuitBreaker.FailureRateThreshold,
			SlowCallThreshold:    30 * time.Second,
			Window:               60 * time.Second,
			OpenDuration:         30 * time.Second,
			HalfOpenMaxRequests:  cfg.CircuitBreaker.HalfOpenMaxRequests,
		}
		executor.breaker = NewCircuitBreaker(breakerCfg)
	}

	if cfg.Routing.Default.Primary != "" {
		routerCfg := RouterConfig{
			DefaultPrimary: cfg.Routing.Default.Primary,
			FallbackList:   cfg.Routing.Default.Fallback,
			MaxAttempts:    2,
			EnableFallback: true,
		}
		executor.router = NewRouter(executor.clientRegistry, routerCfg)
	}

	return executor
}

func (e *Executor) RegisterClient(name string, client LLMClient) {
	e.clientRegistry.Register(name, client)
}

func (e *Executor) SetModelFactory(factory *LLMFactory) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.modelFactory = factory
}

func (e *Executor) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, CallResult) {
	preparedReq, err := e.prepareRequest(req)
	if err != nil {
		return nil, CallResult{
			Provider: req.Provider,
			Model:    req.Model,
			Success:  false,
			Error:    err,
		}
	}
	req = preparedReq

	if err := e.pool.Acquire(ctx, 3*time.Second); err != nil {
		return nil, CallResult{
			Success: false,
			Error:   err,
		}
	}
	defer e.pool.Release()

	if e.breaker != nil && !e.breaker.Allow() {
		return nil, CallResult{
			Success: false,
			Error:   ErrPoolFull,
		}
	}

	start := time.Now()
	resp, result := e.router.Route(ctx, req)
	result.LatencyMs = time.Since(start).Milliseconds()

	if resp != nil && result.Success {
		if e.breaker != nil {
			e.breaker.RecordSuccess(time.Since(start))
		}
		return resp, result
	}

	if e.breaker != nil {
		e.breaker.RecordFailure()
	}
	return resp, result
}

func (e *Executor) Stream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, CallResult) {
	preparedReq, err := e.prepareRequest(req)
	if err != nil {
		return nil, CallResult{
			Provider: req.Provider,
			Model:    req.Model,
			Success:  false,
			Error:    err,
		}
	}
	req = preparedReq

	if err := e.streamPool.Acquire(ctx, 3*time.Second); err != nil {
		return nil, CallResult{
			Success: false,
			Error:   err,
		}
	}

	out := make(chan StreamChunk, 32)

	go func() {
		defer e.streamPool.Release()

		start := time.Now()
		ch, result := e.router.RouteStream(ctx, req)

		if ch == nil {
			out <- StreamChunk{Err: result.Error}
			close(out)
			return
		}

		for chunk := range ch {
			if chunk.Err != nil {
				if e.breaker != nil {
					e.breaker.RecordFailure()
				}
				result.Success = false
				result.LatencyMs = time.Since(start).Milliseconds()
			} else {
				if e.breaker != nil {
					e.breaker.RecordSuccess(time.Since(start))
				}
				result.LatencyMs = time.Since(start).Milliseconds()
				if chunk.Done {
					result.Success = true
				}
			}
			out <- chunk
		}
		close(out)
	}()

	return out, CallResult{Success: true}
}

func (e *Executor) prepareRequest(req ChatRequest) (ChatRequest, error) {
	if req.Provider != "" || req.ModelConfigID == 0 {
		return req, nil
	}

	name := ModelConfigClientName(req.ModelConfigID)
	if _, ok := e.clientRegistry.Get(name); ok {
		req.Provider = name
		return req, nil
	}

	e.mu.RLock()
	factory := e.modelFactory
	e.mu.RUnlock()
	if factory == nil {
		return req, ErrNoProvider
	}

	e.clientRegistry.Register(name, NewEinoModelClient(factory, req.ModelConfigID))
	req.Provider = name
	return req, nil
}

func (e *Executor) Stats() ExecutorStats {
	stats := ExecutorStats{
		Pool:   e.pool.Stats(),
		Stream: e.streamPool.Stats(),
	}

	if e.breaker != nil {
		stats.Breaker = e.breaker.State().String()
	}

	return stats
}

type ExecutorStats struct {
	Pool    PoolStats
	Stream  PoolStats
	Breaker string
}
