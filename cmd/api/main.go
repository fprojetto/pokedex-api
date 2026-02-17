package main

import (
	"context"
	"errors"
	"log"

	"github.com/fprojetto/pokedex-api/application"
	"github.com/fprojetto/pokedex-api/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(errors.Join(errors.New("failed to load configuration for the app"), err))
	}

	if err := application.Run(context.Background(), *cfg); err != nil {
		log.Fatal(err)
	}
}
