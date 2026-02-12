package application

import (
	"context"
	"log"
	"net/http"

	"github.com/fprojetto/pokedex-api/api"
	"github.com/fprojetto/pokedex-api/api/handler"
	"github.com/fprojetto/pokedex-api/config"
	"github.com/fprojetto/pokedex-api/pkg/server"
)

func BuildAPI() http.Handler {
	getPokemonHandler := handler.GetPokemon()
	getPokemonTranslatedHandler := handler.GetPokemonTranslated()
	pokemonMux := api.NewPokemonRouter(
		getPokemonHandler,
		getPokemonTranslatedHandler,
	)

	return pokemonMux
}

func Run(ctx context.Context, cfg config.Config) error {
	// build api
	api := BuildAPI()

	// build and run http server
	httpServer, err := server.NewHTTPServer(server.ServerConfig{
		Addr:            cfg.Addr,
		ShutdownTimeout: cfg.ShutdownTimeout,
		OnShutdown:      shutdown,
	}, api)
	if err != nil {
		return err
	}

	return httpServer.Run(ctx)
}

func shutdown() {
	log.Println("shutting down application")
}
