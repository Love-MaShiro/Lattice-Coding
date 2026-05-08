package chat

import (
	"time"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/redis"
	"lattice-coding/internal/modules/chat/api"
	"lattice-coding/internal/modules/chat/application"
	"lattice-coding/internal/modules/chat/domain"
	"lattice-coding/internal/modules/chat/infra/persistence"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleProvider struct {
	DB           *gorm.DB
	Redis        *redis.Client
	ModelFactory application.ChatModelFactory
	MemoryConfig config.ChatMemoryConfig
}

type Module struct {
	SessionRepo    domain.SessionRepository
	MessageRepo    domain.MessageRepository
	CommandService *application.CommandService
	QueryService   *application.QueryService
	Handler        *api.Handler
}

func NewModule(p *ModuleProvider) *Module {
	_ = persistence.Migrate(p.DB)

	sessionRepo := persistence.NewSessionRepositoryImpl(p.DB)
	messageRepo := persistence.NewMessageRepositoryImpl(p.DB)
	agentGetter := persistence.NewAgentGetter(p.DB)

	cmdSvc := application.NewCommandService(
		sessionRepo,
		messageRepo,
		agentGetter,
		p.ModelFactory,
		p.Redis,
		memoryConfigFromConfig(p.MemoryConfig),
	)
	querySvc := application.NewQueryService(sessionRepo, messageRepo)
	handler := api.NewHandler(cmdSvc, querySvc)

	return &Module{
		SessionRepo:    sessionRepo,
		MessageRepo:    messageRepo,
		CommandService: cmdSvc,
		QueryService:   querySvc,
		Handler:        handler,
	}
}

func RegisterRoutes(group *gin.RouterGroup, m *Module) {
	api.RegisterRoutes(group, m.Handler)
}

func memoryConfigFromConfig(cfg config.ChatMemoryConfig) application.MemoryConfig {
	cacheTTL, err := time.ParseDuration(cfg.CacheTTL)
	if err != nil {
		cacheTTL = 0
	}
	return application.MemoryConfig{
		CompressionThreshold: cfg.CompressionThreshold,
		RetainAfterCompress:  cfg.RetainAfterCompress,
		CacheTTL:             cacheTTL,
	}
}
