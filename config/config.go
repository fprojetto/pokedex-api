package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration

	PokemonAPIURL string
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pokemonAPIURL := os.Getenv("POKEMON_API_URL")
	if pokemonAPIURL == "" {
		return nil, errors.New("missing POKEMON_API_URL environment variable")
	}

	cfg := Config{
		Addr:            ":" + port,
		ShutdownTimeout: 5 * time.Second,

		PokemonAPIURL: pokemonAPIURL,
	}

	return &cfg, nil
}
