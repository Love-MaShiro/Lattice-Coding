package app

import (
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
}

func InitModules(d *Dependencies) *Modules {
	providerModule := provider.NewModule(&provider.ModuleProvider{
		DB: d.MySQL,
	})

	return &Modules{
		ProviderModule: providerModule,
	}
}

func (m *Modules) RegisterRoutes(api *gin.RouterGroup) {
	m.ProviderModule.RegisterRoutes(api)
	agent.RegisterRoutes(api)
	chat.RegisterRoutes(api)
	run.RegisterRoutes(api)
	mcp.RegisterRoutes(api)
	workflow.RegisterRoutes(api)
	knowledge.RegisterRoutes(api)
	safety.RegisterRoutes(api)
	audit.RegisterRoutes(api)
}
