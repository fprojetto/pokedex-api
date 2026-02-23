package application

import (
	"context"
	"log"
	"net/http"

	"github.com/fprojetto/pokedex-api/internal/api"
	"github.com/fprojetto/pokedex-api/internal/api-client/pokeapi"
	"github.com/fprojetto/pokedex-api/internal/api-client/translationapi"
	"github.com/fprojetto/pokedex-api/internal/api/handler"
	"github.com/fprojetto/pokedex-api/internal/config"
	"github.com/fprojetto/pokedex-api/internal/service"
	"github.com/fprojetto/pokedex-api/pkg/server"
)

func BuildAPI(pokemonGetter handler.PokemonGetter, pokemonGetterTranslated handler.PokemonGetterTranslator) http.Handler {
	getPokemonHandler := handler.GetPokemon(pokemonGetter)
	getPokemonTranslatedHandler := handler.GetPokemonTranslated(pokemonGetterTranslated)
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
	translationAPIClient, err := translationapi.NewClient(cfg.TranslationAPIURL)
	if err != nil {
		return err
	}

	pokemonGetterService := service.PokemonGetterService(pokeAPIClient.PokemonInfo)
	pokemonGetterTranslatedService := service.PokemonGetterTranslatorService(
		pokeAPIClient.PokemonInfo,
		translationAPIClient.Translate,
	)
	apiMux := BuildAPI(
		pokemonGetterService,
		pokemonGetterTranslatedService,
	)

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
