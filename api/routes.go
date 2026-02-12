package api

import (
	"net/http"

	"github.com/fprojetto/pokedex-api/pkg/server"
)

func NewPokemonRouter(getPokemon http.HandlerFunc, getPokemonTranslated http.HandlerFunc) http.Handler {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /api/pokemon/{name}", getPokemon)
	apiMux.HandleFunc("GET /api/pokemon/translated/{name}", getPokemonTranslated)

	return server.RequestIDMiddleware(apiMux)
}
