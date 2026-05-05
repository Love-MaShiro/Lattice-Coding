package llm

import (
	"context"
	"sync"
	"time"
)

type SemaphorePool struct {
	sem   chan struct{}
	mu    sync.RWMutex
	stats PoolStats
}

type PoolStats struct {
	Active int64
	Total  int64
	Max    int
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
		p.mu.Lock()
		p.stats.Active++
		p.stats.Total++
		p.mu.Unlock()
		return nil
	case <-timer.C:
		return ErrPoolFull
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *SemaphorePool) Release() {
	p.mu.Lock()
	p.stats.Active--
	p.mu.Unlock()

	select {
	case <-p.sem:
	default:
	}
}

func (p *SemaphorePool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}

func (p *SemaphorePool) Available() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return cap(p.sem) - len(p.sem)
}
