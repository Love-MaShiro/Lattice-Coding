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

	"github.com/gin-gonic/gin"
)

type Modules struct {
	ProviderModule *provider.Module
	AgentModule    *agent.Module
	ChatModule     *chat.Module
}

func InitModules(d *Dependencies) *Modules {
	providerModule := provider.NewModule(&provider.ModuleProvider{
		DB:           d.MySQL,
		Encryptor:    crypto.NewNoopEncryptor(),
		AgentChecker: agent.NewAgentRefCounter(d.MySQL),
	})

	modelConfigChecker := NewProviderModelConfigChecker(providerModule.QueryService)

	agentModule := agent.NewModule(&agent.ModuleProvider{
		DB:                 d.MySQL,
		ModelConfigChecker: modelConfigChecker,
	})

	chatModule := chat.NewModule(&chat.ModuleProvider{
		DB:           d.MySQL,
		Redis:        d.Redis,
		ModelFactory: providerModule.LLMFactory,
		MemoryConfig: d.Config.LLM.ChatMemory,
	})

	return &Modules{
		ProviderModule: providerModule,
		AgentModule:    agentModule,
		ChatModule:     chatModule,
	}
}

func (m *Modules) RegisterRoutes(api *gin.RouterGroup) {
	provider.RegisterRoutes(api, m.ProviderModule)
	agent.RegisterRoutes(api, m.AgentModule)
	chat.RegisterRoutes(api, m.ChatModule)
	run.RegisterRoutes(api)
	mcp.RegisterRoutes(api)
	workflow.RegisterRoutes(api)
	knowledge.RegisterRoutes(api)
	safety.RegisterRoutes(api)
	audit.RegisterRoutes(api)
}
