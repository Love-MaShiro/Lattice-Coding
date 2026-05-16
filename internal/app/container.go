package app

import (
	"gorm.io/gorm"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/common/redis"
	"lattice-coding/internal/runtime/llm"
	runtimetool "lattice-coding/internal/runtime/tool"
)

// Container is the application dependency container. It is intentionally thin
// for now so modules can move away from package-level state incrementally.
type Container struct {
	Config       *config.Config
	Logger       *logger.Logger
	MySQL        *gorm.DB
	Redis        *redis.Client
	VectorDB     *gorm.DB
	LLMExecutor  *llm.Executor
	ToolExecutor *runtimetool.Executor
}

func NewContainer(deps *Dependencies) *Container {
	if deps == nil {
		return &Container{}
	}
	return &Container{
		Config:       deps.Config,
		Logger:       deps.Logger,
		MySQL:        deps.MySQL,
		Redis:        deps.Redis,
		VectorDB:     deps.VectorDB,
		LLMExecutor:  deps.LLMExecutor,
		ToolExecutor: runtimetool.Default(),
	}
}
