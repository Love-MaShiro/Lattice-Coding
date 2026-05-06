package application

import (
	"context"
	"sync"
	"time"

	"lattice-coding/internal/common/errors"
	"lattice-coding/internal/modules/provider/domain"
	"lattice-coding/internal/runtime/llm"
)

type HealthCheckService struct {
	providerRepo       domain.ProviderRepository
	modelConfigRepo    domain.ModelConfigRepository
	providerHealthRepo domain.ProviderHealthRepository
	healthChecker      llm.HealthChecker
	maxConcurrency     int
}

func NewHealthCheckService(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	providerHealthRepo domain.ProviderHealthRepository,
	healthChecker llm.HealthChecker,
) *HealthCheckService {
	return &HealthCheckService{
		providerRepo:       providerRepo,
		modelConfigRepo:    modelConfigRepo,
		providerHealthRepo: providerHealthRepo,
		healthChecker:      healthChecker,
		maxConcurrency:     5,
	}
}

func NewHealthCheckServiceWithConcurrency(
	providerRepo domain.ProviderRepository,
	modelConfigRepo domain.ModelConfigRepository,
	providerHealthRepo domain.ProviderHealthRepository,
	healthChecker llm.HealthChecker,
	maxConcurrency int,
) *HealthCheckService {
	return &HealthCheckService{
		providerRepo:       providerRepo,
		modelConfigRepo:    modelConfigRepo,
		providerHealthRepo: providerHealthRepo,
		healthChecker:      healthChecker,
		maxConcurrency:     maxConcurrency,
	}
}

func (s *HealthCheckService) CheckProvider(ctx context.Context, providerID uint64) (*HealthCheckResult, error) {
	provider, err := s.providerRepo.FindByID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 失败")
	}
	if provider == nil {
		return nil, errors.NotFoundErr("Provider 不存在")
	}

	if !provider.Enabled {
		return &HealthCheckResult{
			ProviderID:   providerID,
			Status:       string(domain.HealthStatusUnknown),
			ErrorCode:    "PROVIDER_DISABLED",
			ErrorMessage: "Provider 已被禁用",
			CheckedAt:    time.Now(),
		}, nil
	}

	modelConfigs, err := s.modelConfigRepo.FindByProviderID(ctx, providerID)
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询模型配置失败")
	}

	var enabledConfig *domain.ModelConfig
	for _, cfg := range modelConfigs {
		if cfg.Enabled {
			enabledConfig = cfg
			break
		}
	}

	if enabledConfig == nil {
		return &HealthCheckResult{
			ProviderID:   providerID,
			Status:       string(domain.HealthStatusUnknown),
			ErrorCode:    "NO_ENABLED_MODEL",
			ErrorMessage: "无可用的模型配置",
			CheckedAt:    time.Now(),
		}, nil
	}

	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.healthChecker.TestModel(checkCtx, enabledConfig.ID)
	if err != nil {
		return nil, errors.InternalWithErr(err, "健康检查失败")
	}

	status := s.determineStatus(result)
	s.updateProviderHealth(ctx, providerID, enabledConfig.ID, status, result)

	return &HealthCheckResult{
		ProviderID:   providerID,
		Status:       status,
		LatencyMs:    result.LatencyMs,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		CheckedAt:    result.CheckedAt,
	}, nil
}

func (s *HealthCheckService) CheckModelConfig(ctx context.Context, modelConfigID uint64) (*HealthCheckResult, error) {
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

	if !provider.Enabled {
		return &HealthCheckResult{
			ProviderID:   provider.ID,
			Status:       string(domain.HealthStatusUnknown),
			ErrorCode:    "PROVIDER_DISABLED",
			ErrorMessage: "Provider 已被禁用",
			CheckedAt:    time.Now(),
		}, nil
	}

	if !modelConfig.Enabled {
		return &HealthCheckResult{
			ProviderID:    provider.ID,
			ModelConfigID: modelConfigID,
			Status:        string(domain.HealthStatusUnknown),
			ErrorCode:     "MODEL_DISABLED",
			ErrorMessage:  "模型配置已被禁用",
			CheckedAt:     time.Now(),
		}, nil
	}

	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.healthChecker.TestModel(checkCtx, modelConfigID)
	if err != nil {
		return nil, errors.InternalWithErr(err, "健康检查失败")
	}

	status := s.determineStatus(result)
	s.updateProviderHealth(ctx, provider.ID, modelConfigID, status, result)

	return &HealthCheckResult{
		ProviderID:    provider.ID,
		ModelConfigID: modelConfigID,
		Status:        status,
		LatencyMs:     result.LatencyMs,
		ErrorCode:     result.ErrorCode,
		ErrorMessage:  result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
	}, nil
}

func (s *HealthCheckService) CheckAllEnabledProviders(ctx context.Context) (*HealthCheckSummary, error) {
	result, err := s.providerRepo.FindPage(ctx, &domain.PageRequest{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, errors.DatabaseErrWithErr(err, "查询 Provider 列表失败")
	}

	var enabledProviders []*domain.Provider
	for _, p := range result.Items {
		if p.Enabled {
			enabledProviders = append(enabledProviders, p)
		}
	}

	summary := &HealthCheckSummary{
		Total: len(enabledProviders),
	}

	if len(enabledProviders) == 0 {
		return summary, nil
	}

	semaphore := make(chan struct{}, s.maxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, provider := range enabledProviders {
		wg.Add(1)
		go func(p *domain.Provider) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			checkCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()

			result, err := s.CheckProvider(checkCtx, p.ID)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				summary.Failed++
				summary.Errors = append(summary.Errors, err.Error())
				return
			}

			switch result.Status {
			case string(domain.HealthStatusHealthy):
				summary.Healthy++
			case string(domain.HealthStatusDegraded):
				summary.Degraded++
			case string(domain.HealthStatusUnhealthy):
				summary.Unhealthy++
			default:
				summary.Failed++
			}
		}(provider)
	}

	wg.Wait()

	return summary, nil
}

func (s *HealthCheckService) determineStatus(result *llm.HealthCheckResult) string {
	if result.Success {
		return string(domain.HealthStatusHealthy)
	}

	switch result.ErrorCode {
	case "TIMEOUT":
		return string(domain.HealthStatusDegraded)
	case "RATE_LIMIT", "429":
		return string(domain.HealthStatusDegraded)
	case "AUTH_ERROR", "UNAUTHORIZED", "401":
		return string(domain.HealthStatusUnhealthy)
	default:
		return string(domain.HealthStatusUnhealthy)
	}
}

func (s *HealthCheckService) updateProviderHealth(ctx context.Context, providerID uint64, modelConfigID uint64, status string, result *llm.HealthCheckResult) {
	s.providerRepo.UpdateHealthStatus(ctx, providerID, domain.HealthStatus(status), result.ErrorMessage)

	health := &domain.ProviderHealth{
		ProviderID:    providerID,
		ModelConfigID: modelConfigID,
		Status:        status,
		LatencyMs:     result.LatencyMs,
		ErrorCode:     result.ErrorCode,
		ErrorMessage:  result.ErrorMessage,
		CheckedAt:     result.CheckedAt,
		CreatedAt:     time.Now(),
	}
	s.providerHealthRepo.Create(ctx, health)
}

type HealthCheckResult struct {
	ProviderID    uint64
	ModelConfigID uint64
	Status        string
	LatencyMs     int64
	ErrorCode     string
	ErrorMessage  string
	CheckedAt     time.Time
}

type HealthCheckSummary struct {
	Total     int
	Healthy   int
	Degraded  int
	Unhealthy int
	Failed    int
	Errors    []string
}
