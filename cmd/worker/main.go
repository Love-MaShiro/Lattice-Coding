package main

import (
	"lattice-coding/internal/common/config"
	"lattice-coding/internal/common/logger"
	"lattice-coding/internal/modules/run"
	"lattice-coding/internal/runtime/eino"
	"lattice-coding/internal/runtime/event"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg)

	eino.Init(cfg)
	event.Init(cfg)

	run.StartWorker(cfg, log)

	log.Info("Worker started successfully")
	select {}
}
