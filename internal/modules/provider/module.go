package provider

import (
	"lattice-coding/internal/modules/provider/api"
	"lattice-coding/internal/modules/provider/application"
	"lattice-coding/internal/modules/provider/infra/persistence"
	"lattice-coding/internal/runtime/llm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModuleProvider 模块依赖
type ModuleProvider struct {
	DB *gorm.DB
}

// Module 模块
type Module struct {
	ProviderHandler *api.ProviderHandler
	LLMFactory      *llm.LLMFactory
}

// NewModule 创建模块
func NewModule(p *ModuleProvider) *Module {
	// 自动迁移表结构
	_ = persistence.Migrate(p.DB)

	providerRepo := persistence.NewProviderRepositoryImpl(p.DB)
	modelConfigRepo := persistence.NewModelConfigRepositoryImpl(p.DB)
	cryptoService := persistence.NewCryptoService()
	providerService := application.NewProviderService(providerRepo, modelConfigRepo, cryptoService)
	llmFactory := llm.NewLLMFactory(providerRepo, modelConfigRepo, cryptoService)
	providerHandler := api.NewProviderHandler(providerService, llmFactory)

	return &Module{
		ProviderHandler: providerHandler,
		LLMFactory:      llmFactory,
	}
}

// RegisterRoutes 注册路由
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	// Providers
	r := api.Group("/v1/providers")
	{
		r.GET("", m.ProviderHandler.ListProviders)
		r.GET("/:id", m.ProviderHandler.GetProvider)
		r.POST("", m.ProviderHandler.CreateProvider)
		r.PUT("/:id", m.ProviderHandler.UpdateProvider)
		r.DELETE("/:id", m.ProviderHandler.DeleteProvider)
		r.POST("/:id/enable", m.ProviderHandler.EnableProvider)
		r.POST("/:id/disable", m.ProviderHandler.DisableProvider)
		r.POST("/:id/test", m.ProviderHandler.TestProvider)
	}

	// Model Configs
	mc := api.Group("/v1/model-configs")
	{
		mc.GET("", m.ProviderHandler.ListModelConfigs)
		mc.POST("", m.ProviderHandler.CreateModelConfig)
		mc.POST("/:id/test", m.ProviderHandler.TestModelConfig)
	}
}
