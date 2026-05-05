package app

import (
	"errors"

	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/runtime/eino"
	"lattice-coding/internal/runtime/event"
	"lattice-coding/internal/runtime/tool"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Bootstrap struct {
	Deps    *Dependencies
	Modules *Modules
	Engine  *gin.Engine
}

func NewHTTPBootstrap() (*Bootstrap, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	log := logger.NewLogger(cfg)

	deps, err := NewDependencies(cfg, log)
	if err != nil {
		return nil, err
	}

	eino.Init(cfg)
	tool.Init(cfg)
	event.Init(cfg)

	engine, api := NewRouter(deps)

	modules := InitModules(deps)
	modules.RegisterRoutes(api)

	return &Bootstrap{
		Deps:    deps,
		Modules: modules,
		Engine:  engine,
	}, nil
}

func NewWorkerBootstrap() (*Bootstrap, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	log := logger.NewLogger(cfg)

	deps, err := NewDependencies(cfg, log)
	if err != nil {
		return nil, err
	}

	eino.Init(cfg)
	tool.Init(cfg)
	event.Init(cfg)

	modules := InitModules(deps)

	return &Bootstrap{
		Deps:    deps,
		Modules: modules,
	}, nil
}

func (b *Bootstrap) Close() error {
	if b == nil || b.Deps == nil {
		return nil
	}
	return b.Deps.Close()
}

func (b *Bootstrap) RunHTTP() error {
	if b == nil || b.Deps == nil || b.Engine == nil {
		return errors.New("http bootstrap not initialized")
	}

	b.Deps.Logger.Info("API server starting", zap.String("port", b.Deps.Config.App.Port))
	return b.Engine.Run(":" + b.Deps.Config.App.Port)
}
