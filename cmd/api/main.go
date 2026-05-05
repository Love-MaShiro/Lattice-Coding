package main

import (
	"log"

	"lattice-coding/internal/app"
)

func main() {
	bootstrap, err := app.NewHTTPBootstrap()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = bootstrap.Close() }()

	if err := bootstrap.RunHTTP(); err != nil {
		log.Fatal(err)
	}
}
