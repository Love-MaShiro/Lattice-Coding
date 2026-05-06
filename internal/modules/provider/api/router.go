package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	providers := api.Group("/v1/providers")
	{
		providers.POST("", h.CreateProvider)
		providers.GET("", h.ListProviders)
		providers.GET("/:id", h.GetProvider)
		providers.PUT("/:id", h.UpdateProvider)
		providers.DELETE("/:id", h.DeleteProvider)
		providers.POST("/:id/enable", h.EnableProvider)
		providers.POST("/:id/disable", h.DisableProvider)
		providers.POST("/:id/test", h.TestProvider)
		providers.POST("/:id/sync-models", h.SyncProviderModels)
		providers.POST("/:id/health-check", h.HealthCheckProvider)
		providers.GET("/:id/health", h.GetProviderHealth)
	}

	modelConfigs := api.Group("/v1/model-configs")
	{
		modelConfigs.POST("", h.CreateModelConfig)
		modelConfigs.GET("", h.ListModelConfigs)
		modelConfigs.GET("/:id", h.GetModelConfig)
		modelConfigs.PUT("/:id", h.UpdateModelConfig)
		modelConfigs.DELETE("/:id", h.DeleteModelConfig)
		modelConfigs.POST("/:id/enable", h.EnableModelConfig)
		modelConfigs.POST("/:id/disable", h.DisableModelConfig)
		modelConfigs.POST("/:id/test", h.TestModelConfig)
		modelConfigs.POST("/:id/default", h.SetDefaultModelConfig)
		modelConfigs.GET("/:id/health", h.GetModelConfigHealth)
	}
}
