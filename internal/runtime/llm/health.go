package llm

import (
	"context"
	"time"
)

type HealthCheckResult struct {
	Success      bool
	ProviderID   uint64
	ModelConfigID uint64
	LatencyMs    int64
	HealthStatus string
	ErrorCode    string
	ErrorMessage string
	CheckedAt    time.Time
}

type HealthChecker interface {
	TestModel(ctx context.Context, modelConfigID uint64) (*HealthCheckResult, error)
	TestProvider(ctx context.Context, providerID uint64) (*HealthCheckResult, error)
}
