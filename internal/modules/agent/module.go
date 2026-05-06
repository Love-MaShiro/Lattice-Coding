package agent

import (
	"lattice-coding/internal/modules/agent/api"
	"lattice-coding/internal/modules/agent/application"
	"lattice-coding/internal/modules/agent/domain"
	"lattice-coding/internal/modules/agent/infra/persistence"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB *gorm.DB
}

type Module struct {
	AgentRepo      domain.AgentRepository
	CommandService *application.CommandService
	QueryService   *application.QueryService
	Handler        *api.Handler
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	agentRepo := persistence.NewAgentRepositoryImpl(p.DB)
	providerGetter := persistence.NewProviderGetter(p.DB)

	cmdSvc := application.NewCommandService(agentRepo, providerGetter)
	querySvc := application.NewQueryService(agentRepo)
	handler := api.NewHandler(cmdSvc, querySvc)

	return &Module{
		AgentRepo:      agentRepo,
		CommandService: cmdSvc,
		QueryService:   querySvc,
		Handler:        handler,
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	api.RegisterRoutes(group, m.Handler)
}

func NewAgentRefCounter(db *gorm.DB) domain.AgentReferenceChecker {
	return persistence.NewAgentRefCounter(db)
}
