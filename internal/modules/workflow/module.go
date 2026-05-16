package workflow

import (
	rundomain "lattice-coding/internal/modules/run/domain"
	"lattice-coding/internal/modules/workflow/api"
	"lattice-coding/internal/modules/workflow/application"
	"lattice-coding/internal/modules/workflow/domain"
	"lattice-coding/internal/modules/workflow/infra/persistence"
	"lattice-coding/internal/runtime/llm"
	runtimetool "lattice-coding/internal/runtime/tool"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB            *gorm.DB
	RunRepo       rundomain.RunRepository
	LLMExecutor   *llm.Executor
	ToolExecutor  *runtimetool.ToolExecutor
	AuditRecorder runtimetool.AuditRecorder
}

type Module struct {
	WorkflowRepo  domain.WorkflowRepository
	NodeRegistry  *application.NodeRegistry
	CommandSvc    *application.CommandService
	QuerySvc      *application.QueryService
	CodeReviewSvc *application.CodeReviewService
	Handler       *api.Handler
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	parser := domain.NewJSONNodeConfigParser()
	workflowRepo := persistence.NewWorkflowRepositoryImpl(p.DB, parser)
	nodeRegistry := application.NewNodeRegistry()
	commandSvc := application.NewCommandService(workflowRepo, parser)
	querySvc := application.NewQueryService(workflowRepo)
	codeReviewSvc := application.NewCodeReviewService(p.RunRepo, p.ToolExecutor, p.LLMExecutor, p.AuditRecorder)
	handler := api.NewHandler(commandSvc, querySvc, codeReviewSvc)

	return &Module{
		WorkflowRepo:  workflowRepo,
		NodeRegistry:  nodeRegistry,
		CommandSvc:    commandSvc,
		QuerySvc:      querySvc,
		CodeReviewSvc: codeReviewSvc,
		Handler:       handler,
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	api.RegisterRoutes(group, m.Handler)
}
