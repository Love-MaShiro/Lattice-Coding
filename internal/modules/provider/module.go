package provider

import (
	"lattice-coding/internal/common/crypto"
	"lattice-coding/internal/modules/provider/api"
	"lattice-coding/internal/modules/provider/application"
	"lattice-coding/internal/modules/provider/domain"
	"lattice-coding/internal/modules/provider/infra/persistence"
	"lattice-coding/internal/runtime/llm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB           *gorm.DB
	Encryptor    crypto.Encryptor
	AgentChecker application.AgentReferenceChecker
}

type Module struct {
	ProviderRepo       domain.ProviderRepository
	ModelConfigRepo    domain.ModelConfigRepository
	ProviderHealthRepo domain.ProviderHealthRepository
	CommandService     *application.CommandService
	QueryService       *application.QueryService
	HealthService      *application.HealthService
	HealthCheckService *application.HealthCheckService
	SyncService        *application.SyncService
	LLMFactory         *llm.LLMFactory
	Handler            *api.Handler
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	encryptor := p.Encryptor
	if encryptor == nil {
		encryptor = crypto.NewNoopEncryptor()
	}

	providerRepo := persistence.NewProviderRepositoryImpl(p.DB)
	modelConfigRepo := persistence.NewModelConfigRepositoryImpl(p.DB)
	healthRepo := persistence.NewProviderHealthRepositoryImpl(p.DB)

	llmFactory := llm.NewLLMFactory(providerRepo, healthRepo, modelConfigRepo)
	modelLister := llm.NewModelLister()

	cmdSvc := application.NewCommandService(
		providerRepo,
		modelConfigRepo,
		encryptor,
		p.AgentChecker,
	)
	querySvc := application.NewQueryService(
		providerRepo,
		modelConfigRepo,
	)
	healthSvc := application.NewHealthService(
		providerRepo,
		modelConfigRepo,
		healthRepo,
		llmFactory,
	)
	healthCheckSvc := application.NewHealthCheckService(
		providerRepo,
		modelConfigRepo,
		healthRepo,
		llmFactory,
	)
	syncSvc := application.NewSyncService(
		providerRepo,
		modelConfigRepo,
		modelLister,
	)
	handler := api.NewHandler(
		cmdSvc,
		querySvc,
		healthSvc,
		healthCheckSvc,
		syncSvc,
	)

	return &Module{
		ProviderRepo:       providerRepo,
		ModelConfigRepo:    modelConfigRepo,
		ProviderHealthRepo: healthRepo,
		CommandService:     cmdSvc,
		QueryService:       querySvc,
		HealthService:      healthSvc,
		HealthCheckService: healthCheckSvc,
		SyncService:        syncSvc,
		Handler:            handler,
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	api.RegisterRoutes(group, m.Handler)
}
