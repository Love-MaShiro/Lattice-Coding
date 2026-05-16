package run

import (
	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/modules/run/api"
	"lattice-coding/internal/modules/run/application"
	"lattice-coding/internal/modules/run/domain"
	"lattice-coding/internal/modules/run/infra/persistence"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB *gorm.DB
}

type Module struct {
	RunRepo            domain.RunRepository
	ToolInvocationRepo domain.ToolInvocationRepository
	QuerySvc           *application.QueryService
	InvocationRecorder *application.ToolInvocationRecorder
	Handler            *api.Handler
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	runRepo := persistence.NewRunRepositoryImpl(p.DB)
	invocationRepo := persistence.NewToolInvocationRepositoryImpl(p.DB)
	querySvc := application.NewQueryService(runRepo, invocationRepo)
	invocationRecorder := application.NewToolInvocationRecorder(invocationRepo)
	handler := api.NewHandler(querySvc)

	return &Module{
		RunRepo:            runRepo,
		ToolInvocationRepo: invocationRepo,
		QuerySvc:           querySvc,
		InvocationRecorder: invocationRecorder,
		Handler:            handler,
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	api.RegisterRoutes(group, m.Handler)
}

func StartWorker(cfg *config.Config, log *logger.Logger) {
	log.Info("Starting agent run worker", zap.String("queue", cfg.Worker.QueueName))
}
