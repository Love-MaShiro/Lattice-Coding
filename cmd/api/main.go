package main

import (
	"lattice-coding/internal/app"
	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/common/middleware"
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
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg)

	bootstrap, err := app.NewBootstrap(cfg)
	if err != nil {
		log.Warn("failed to initialize bootstrap, running without DB/Redis", zap.Error(err))
	} else {
		defer func() {
			if err := bootstrap.Close(); err != nil {
				log.Error("failed to close bootstrap", zap.Error(err))
			}
		}()
	}

	r := gin.New()

	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))
	r.Use(middleware.CORS())
	r.Use(middleware.Trace())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    "Lattice-coding is running",
		})
	})

	provider.RegisterRoutes(r)
	agent.RegisterRoutes(r)
	chat.RegisterRoutes(r)
	run.RegisterRoutes(r)
	mcp.RegisterRoutes(r)
	workflow.RegisterRoutes(r)
	knowledge.RegisterRoutes(r)
	safety.RegisterRoutes(r)
	audit.RegisterRoutes(r)

	log.Info("API server starting", zap.String("port", cfg.App.Port))
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}
