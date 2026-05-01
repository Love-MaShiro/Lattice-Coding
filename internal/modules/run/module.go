package run

import (
	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"

	"go.uber.org/zap"
)

func RegisterRoutes(r interface{}) {
}

func StartWorker(cfg *config.Config, log *logger.Logger) {
	log.Info("Starting agent run worker", zap.String("queue", cfg.Worker.QueueName))
}
