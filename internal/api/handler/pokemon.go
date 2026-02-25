package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/fprojetto/pokedex-api/internal/api"
	"github.com/fprojetto/pokedex-api/internal/model"
	"github.com/fprojetto/pokedex-api/internal/service"
)

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary *bool  `json:"isLegendary"`
}

type PokemonGetter func(ctx context.Context, name string) (model.Pokemon, error)
type PokemonGetterTranslator func(ctx context.Context, name string) (model.Pokemon, error)

func GetPokemon(getPokemon PokemonGetter) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")
		if name == "" {
			api.WriteError(w, req, http.StatusBadRequest, api.ErrCodeBadRequest, "missing name parameter")
			return
		}
		p, err := getPokemon(req.Context(), name)
		if err != nil {
			handleError(w, req, err)
			return
		}

		pokemon := mapper(p)

		api.WriteJSON(w, req, pokemon, http.StatusOK)
	}
}

func GetPokemonTranslated(
	getPokemonTranslated PokemonGetterTranslator,
) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")
		if name == "" {
			api.WriteError(w, req, http.StatusBadRequest, api.ErrCodeBadRequest, "missing name parameter")
			return
		}

		p, err := getPokemonTranslated(req.Context(), name)
		if err != nil {
			handleError(w, req, err)
			return
		}

		pokemon := mapper(p)

		api.WriteJSON(w, req, pokemon, http.StatusOK)
	}
}

func mapper(p model.Pokemon) Pokemon {
	return Pokemon{
		Name:        p.Name,
		Description: p.Description,
		Habitat:     p.Habitat,
		IsLegendary: p.IsLegendary,
	}
}

func handleError(w http.ResponseWriter, req *http.Request, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		log.Printf("not found error: %v", err)
		api.WriteError(w, req, http.StatusNotFound, api.ErrCodeNotFound, "resource not found")
	default:
		// Log the actual error for internal tracking, but return a generic message to the client
		log.Printf("internal server error: %v", err)
		api.WriteError(w, req, http.StatusInternalServerError, api.ErrCodeInternal, "internal server error")
	}
}
