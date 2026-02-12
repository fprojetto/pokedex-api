package config

import (
	"os"
	"time"
)

type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := Config{
		Addr:            ":" + port,
		ShutdownTimeout: 5 * time.Second,
	}

	return cfg, nil
}
