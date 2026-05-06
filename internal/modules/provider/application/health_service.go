package application

import (
	"context"
	"time"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
	"lattice-coding/internal/runtime/llm"
)

type HealthService struct {
	providerRepo       domain.ProviderRepository
	modelConfigRepo    domain.ModelConfigRepository
	providerHealthRepo domain.ProviderHealthRepository
	healthChecker      llm.HealthChecker
}

func NewHealthService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	providerHealthRepo domain.ProviderHealthRepository,
	healthChecker llm.HealthChecker,
) *HealthService {
	return &HealthService{
		providerRepo:       providerRepo,
		modelConfigRepo:    modelConfigRepo,
		providerHealthRepo: providerHealthRepo,
		healthChecker:      healthChecker,
	}
}

func (s *HealthService) TestProvider(ctx context.Context, providerID uint64) (*HealthTestResult, error) {
	provider, err := s.providerRepo.FindByID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	result, err := s.healthChecker.TestProvider(ctx, providerID)
	if err != nil {
		return nil, errors.InternalWithErr(err, "健康检查失败")
	}

	s.updateProviderHealth(ctx, providerID, 0, result)

	return &HealthTestResult{
		Success:       result.Success,
		ProviderID:    result.ProviderID,
		LatencyMs:     result.LatencyMs,
		HealthStatus:  result.HealthStatus,
		ErrorCode:     result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
	}, nil
}

func (s *HealthService) TestModelConfig(ctx context.Context, modelConfigID uint64) (*ModelConfigTestResult, error) {
	modelConfig, err := s.modelConfigRepo.FindByID(ctx, modelConfigID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}
	if modelConfig == nil {
		return nil, errors.NotFoundErr("模型配置不存在")
	}

	provider, err := s.providerRepo.FindByID(ctx, modelConfig.ProviderID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	result, err := s.healthChecker.TestModel(ctx, modelConfigID)
	if err != nil {
		return nil, errors.InternalWithErr(err, "健康检查失败")
	}

	s.updateProviderHealth(ctx, provider.ID, modelConfigID, result)

	return &ModelConfigTestResult{
		Success:       result.Success,
		ProviderID:    result.ProviderID,
		ModelConfigID: result.ModelConfigID,
		LatencyMs:     result.LatencyMs,
		HealthStatus:  result.HealthStatus,
		ErrorCode:     result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
	}, nil
}

func (s *HealthService) updateProviderHealth(ctx context.Context, providerID uint64, modelConfigID uint64, result *llm.HealthCheckResult) {
	s.providerRepo.UpdateHealthStatus(ctx, providerID, domain.HealthStatus(result.HealthStatus), result.ErrorMessage)

	health := &domain.ProviderHealth{
		ProviderID:    providerID,
		ModelConfigID: modelConfigID,
		Status:        result.HealthStatus,
		LatencyMs:     result.LatencyMs,
		ErrorCode:     result.ErrorCode,
		ErrorMessage:  result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
		CreatedAt:     time.Now(),
	}
	s.providerHealthRepo.Create(ctx, health)
}

type HealthTestResult struct {
	Success      bool
	ProviderID   uint64
	LatencyMs    int64
	HealthStatus string
	ErrorCode    string
	ErrorMessage string
	CheckedAt    time.Time
}

type ModelConfigTestResult struct {
	Success       bool
	ProviderID    uint64
	ModelConfigID uint64
	LatencyMs     int64
	HealthStatus  string
	ErrorCode     string
	ErrorMessage  string
	CheckedAt     time.Time
}
