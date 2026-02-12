package api

import (
	"net/http"
)

func NewPokemonRouter(getPokemon http.HandlerFunc, getPokemonTranslated http.HandlerFunc) *http.ServeMux {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /api/pokemon/{name}", getPokemon)
	apiMux.HandleFunc("GET /api/pokemon/translated/{name}", getPokemonTranslated)

	return apiMux
}
