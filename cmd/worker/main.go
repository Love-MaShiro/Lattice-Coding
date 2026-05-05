package main

import (
	"log"

	"lattice-coding/internal/app"
	"lattice-coding/internal/modules/run"
)

func main() {
	bootstrap, err := app.NewWorkerBootstrap()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = bootstrap.Close() }()

	run.StartWorker(bootstrap.Deps.Config, bootstrap.Deps.Logger)

	bootstrap.Deps.Logger.Info("Worker started successfully")
	select {}
}
