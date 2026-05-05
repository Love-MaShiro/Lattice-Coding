package run

import (
	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(api *gin.RouterGroup) {
	r := api.Group("/v1/runs")
	{
		r.GET("", listRuns)
		r.POST("", createRun)
		r.GET("/:run_id/events", getRunEvents)
	}
}

func StartWorker(cfg *config.Config, log *logger.Logger) {
	log.Info("Starting agent run worker", zap.String("queue", cfg.Worker.QueueName))
}

func listRuns(c *gin.Context) {}

func createRun(c *gin.Context) {}

func getRunEvents(c *gin.Context) {}
