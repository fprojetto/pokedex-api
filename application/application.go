package application

import (
	"context"
	"log"
	"net/http"

	"github.com/fprojetto/pokedex-api/api"
	"github.com/fprojetto/pokedex-api/api-client/pokeapi"
	"github.com/fprojetto/pokedex-api/api/handler"
	"github.com/fprojetto/pokedex-api/config"
	"github.com/fprojetto/pokedex-api/pkg/server"
	"github.com/fprojetto/pokedex-api/service"
)

func BuildAPI(pokemonGetter handler.PokemonGetter) http.Handler {
	getPokemonHandler := handler.GetPokemon(pokemonGetter)
	getPokemonTranslatedHandler := handler.GetPokemonTranslated()
	pokemonMux := api.NewPokemonRouter(
		getPokemonHandler,
		getPokemonTranslatedHandler,
	)

	return pokemonMux
}

func Run(ctx context.Context, cfg config.Config) error {
	// build api
	pokeAPIClient, err := pokeapi.NewClient(cfg.PokemonAPIURL)
	if err != nil {
		return err
	}
	pokemonGetterService := service.PokemonGetterService(pokeAPIClient.PokemonInfo)
	apiMux := BuildAPI(pokemonGetterService)

	// build and run http server
	httpServer, err := server.NewHTTPServer(server.ServerConfig{
		Addr:            cfg.Addr,
		ShutdownTimeout: cfg.ShutdownTimeout,
		OnShutdown:      shutdown,
	}, apiMux)
	if err != nil {
		return err
	}

	return httpServer.Run(ctx)
}

func shutdown() {
	log.Println("shutting down application")
}
