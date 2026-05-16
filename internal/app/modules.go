package app

import (
	"lattice-coding/internal/common/crypto"
	"lattice-coding/internal/modules/agent"
	"lattice-coding/internal/modules/audit"
	"lattice-coding/internal/modules/chat"
	"lattice-coding/internal/modules/knowledge"
	"lattice-coding/internal/modules/mcp"
	"lattice-coding/internal/modules/provider"
	"lattice-coding/internal/modules/run"
	"lattice-coding/internal/modules/safety"
	"lattice-coding/internal/modules/workflow"
	runtimetool "lattice-coding/internal/runtime/tool"
	"lattice-coding/internal/runtime/tool/builtin"

	"github.com/gin-gonic/gin"
)

type Modules struct {
	ProviderModule  *provider.Module
	AgentModule     *agent.Module
	ChatModule      *chat.Module
	KnowledgeModule *knowledge.Module
	WorkflowModule  *workflow.Module
	RunModule       *run.Module
	AuditModule     *audit.Module
}

func InitModules(d *Dependencies) *Modules {
	providerModule := provider.NewModule(&provider.ModuleProvider{
		DB:           d.MySQL,
		Encryptor:    crypto.NewNoopEncryptor(),
		AgentChecker: agent.NewAgentRefCounter(d.MySQL),
	})
	d.LLMExecutor.SetModelFactory(providerModule.LLMFactory)

	modelConfigChecker := NewProviderModelConfigChecker(providerModule.QueryService)

	agentModule := agent.NewModule(&agent.ModuleProvider{
		DB:                 d.MySQL,
		ModelConfigChecker: modelConfigChecker,
	})

	chatModule := chat.NewModule(&chat.ModuleProvider{
		DB:           d.MySQL,
		Redis:        d.Redis,
		LLMExecutor:  d.LLMExecutor,
		MemoryConfig: d.Config.LLM.ChatMemory,
	})
	knowledgeModule := knowledge.NewModule(&knowledge.ModuleProvider{
		DB: d.MySQL,
	})
	runModule := run.NewModule(&run.ModuleProvider{
		DB: d.MySQL,
	})
	auditModule := audit.NewModule(&audit.ModuleProvider{
		DB: d.MySQL,
	})
	readStateManager := runtimetool.NewInMemoryFileReadStateManager()
	_ = builtin.RegisterCodingTools(runtimetool.Default().Registry(), readStateManager)
	runtimetool.Default().SetToolInvocationRecorder(runModule.InvocationRecorder)
	runtimetool.Default().SetAuditRecorder(auditModule.AuditRecorder)
	workflowModule := workflow.NewModule(&workflow.ModuleProvider{
		DB:            d.MySQL,
		RunRepo:       runModule.RunRepo,
		LLMExecutor:   d.LLMExecutor,
		ToolExecutor:  runtimetool.Default(),
		AuditRecorder: auditModule.AuditRecorder,
	})

	return &Modules{
		ProviderModule:  providerModule,
		AgentModule:     agentModule,
		ChatModule:      chatModule,
		KnowledgeModule: knowledgeModule,
		WorkflowModule:  workflowModule,
		RunModule:       runModule,
		AuditModule:     auditModule,
	}
}

func (m *Modules) RegisterRoutes(api *gin.RouterGroup) {
	provider.RegisterRoutes(api, m.ProviderModule)
	agent.RegisterRoutes(api, m.AgentModule)
	chat.RegisterRoutes(api, m.ChatModule)
	run.RegisterRoutes(api, m.RunModule)
	mcp.RegisterRoutes(api)
	workflow.RegisterRoutes(api, m.WorkflowModule)
	knowledge.RegisterRoutes(api, m.KnowledgeModule)
	safety.RegisterRoutes(api)
	audit.RegisterRoutes(api, m.AuditModule)
}
