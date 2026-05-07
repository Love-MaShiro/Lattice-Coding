package app

import (
	"context"

	"lattice-coding/internal/modules/agent/application"
	providerApp "lattice-coding/internal/modules/provider/application"
)

type ProviderModelConfigChecker struct {
	providerQuerySvc *providerApp.QueryService
}

func NewProviderModelConfigChecker(providerQuerySvc *providerApp.QueryService) application.ModelConfigChecker {
	return &ProviderModelConfigChecker{
		providerQuerySvc: providerQuerySvc,
	}
}

func (c *ProviderModelConfigChecker) CheckModelConfigEnabled(ctx context.Context, modelConfigID uint64) error {
	_, err := c.providerQuerySvc.GetModelConfigForRuntime(ctx, modelConfigID)
	return err
}
