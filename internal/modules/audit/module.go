package audit

import (
	"lattice-coding/internal/modules/audit/application"
	"lattice-coding/internal/modules/audit/domain"
	"lattice-coding/internal/modules/audit/infra/persistence"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB *gorm.DB
}

type Module struct {
	AuditLogRepo  domain.AuditLogRepository
	AuditRecorder *application.AuditRecorder
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	auditRepo := persistence.NewAuditLogRepositoryImpl(p.DB)
	auditRecorder := application.NewAuditRecorder(auditRepo)
	return &Module{
		AuditLogRepo:  auditRepo,
		AuditRecorder: auditRecorder,
	}
}

func RegisterRoutes(api *gin.RouterGroup, _ *Module) {
	r := api.Group("/v1/audit")
	{
		r.GET("/logs", listLogs)
		r.GET("/logs/:id", getLog)
		r.GET("/llm-calls", listLLMCalls)
		r.GET("/tool-calls", listToolCalls)
	}
}

func listLogs(c *gin.Context) {}

func getLog(c *gin.Context) {}

func listLLMCalls(c *gin.Context) {}

func listToolCalls(c *gin.Context) {}
