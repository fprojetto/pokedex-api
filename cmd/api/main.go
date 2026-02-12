package main

import (
	"context"
	"log"

	"github.com/fprojetto/pokedex-api/application"
	"github.com/fprojetto/pokedex-api/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(context.Background(), cfg); err != nil {
		log.Fatal(err)
	}
}
