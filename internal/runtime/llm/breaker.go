package llm

import (
	"sync"
	"sync/atomic"
	"time"
)

type CircuitState int32

const (
	StateClosed   CircuitState = 0
	StateOpen     CircuitState = 1
	StateHalfOpen CircuitState = 2
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

type CircuitBreaker struct {
	name                 string
	failureRateThreshold float64
	slowCallThreshold    time.Duration
	window               time.Duration
	openDuration         time.Duration
	halfOpenMaxRequests  int32

	mu              sync.RWMutex
	state           CircuitState
	failures        int32
	successes       int32
	slowCalls       int32
	lastFailure     time.Time
	lastStateChange time.Time
}

type BreakerConfig struct {
	Name                 string
	FailureRateThreshold float64
	SlowCallThreshold    time.Duration
	Window               time.Duration
	OpenDuration         time.Duration
	HalfOpenMaxRequests  int
}

func NewCircuitBreaker(cfg BreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		name:                 cfg.Name,
		failureRateThreshold: cfg.FailureRateThreshold,
		slowCallThreshold:    cfg.SlowCallThreshold,
		window:               cfg.Window,
		openDuration:         cfg.OpenDuration,
		halfOpenMaxRequests:  int32(cfg.HalfOpenMaxRequests),
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	return CircuitState(atomic.LoadInt32((*int32)(&cb.state)))
}

func (cb *CircuitBreaker) Allow() bool {
	switch cb.State() {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateChange) > cb.openDuration {
			cb.halfOpen()
			return true
		}
		return false
	case StateHalfOpen:
		return atomic.LoadInt32(&cb.successes) < cb.halfOpenMaxRequests
	default:
		return false
	}
}

func (cb *CircuitBreaker) RecordSuccess(duration time.Duration) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.State() == StateHalfOpen {
		atomic.AddInt32(&cb.successes, 1)
		if atomic.LoadInt32(&cb.successes) >= cb.halfOpenMaxRequests {
			cb.close()
		}
		return
	}

	atomic.AddInt32(&cb.successes, 1)
	if duration > cb.slowCallThreshold {
		atomic.AddInt32(&cb.slowCalls, 1)
	}

	cb.reset()
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	atomic.AddInt32(&cb.failures, 1)
	cb.lastFailure = time.Now()

	if cb.State() == StateHalfOpen {
		cb.trip()
		return
	}

	total := atomic.LoadInt32(&cb.successes) + atomic.LoadInt32(&cb.failures)
	if total >= 10 && cb.failureRate() >= cb.failureRateThreshold {
		cb.trip()
	}
}

func (cb *CircuitBreaker) close() {
	atomic.StoreInt32((*int32)(&cb.state), int32(StateClosed))
	cb.lastStateChange = time.Now()
	atomic.StoreInt32(&cb.failures, 0)
	atomic.StoreInt32(&cb.successes, 0)
	atomic.StoreInt32(&cb.slowCalls, 0)
}

func (cb *CircuitBreaker) reset() {
	atomic.StoreInt32(&cb.failures, 0)
	atomic.StoreInt32(&cb.successes, 0)
	atomic.StoreInt32(&cb.slowCalls, 0)
}

func (cb *CircuitBreaker) trip() {
	atomic.StoreInt32((*int32)(&cb.state), int32(StateOpen))
	cb.lastStateChange = time.Now()
}

func (cb *CircuitBreaker) halfOpen() {
	atomic.StoreInt32((*int32)(&cb.state), int32(StateHalfOpen))
	cb.lastStateChange = time.Now()
	atomic.StoreInt32(&cb.successes, 0)
	atomic.StoreInt32(&cb.failures, 0)
	atomic.StoreInt32(&cb.slowCalls, 0)
}

func (cb *CircuitBreaker) failureRate() float64 {
	total := atomic.LoadInt32(&cb.successes) + atomic.LoadInt32(&cb.failures)
	if total == 0 {
		return 0
	}
	return float64(atomic.LoadInt32(&cb.failures)) / float64(total)
}
