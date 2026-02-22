package handler

import (
	"context"
	"net/http"

	"github.com/fprojetto/pokedex-api/api"
	"github.com/fprojetto/pokedex-api/model"
	"github.com/fprojetto/pokedex-api/service"
)

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary *bool  `json:"isLegendary"`
}

type PokemonGetter func(ctx context.Context, name string) (model.Pokemon, error)

func GetPokemon(getPokemon PokemonGetter) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")
		if name == "" {
			api.WriteError(w, req, http.StatusBadRequest, api.ErrCodeBadRequest, "missing name parameter")
			return
		}
		p, err := getPokemon(req.Context(), name)
		if err != nil {
			switch err {
			case service.ErrNotFound:
				api.WriteError(w, req, http.StatusNotFound, api.ErrCodeNotFound, err.Error())
			default:
				api.WriteError(w, req, http.StatusInternalServerError, api.ErrCodeInternal, err.Error())
			}
			return
		}

		pokemon := Pokemon{
			Name:        p.Name,
			Description: p.Description,
			Habitat:     p.Habitat,
			IsLegendary: p.IsLegendary,
		}

		api.WriteJSON(w, req, pokemon, http.StatusOK)
	}
}

func GetPokemonTranslated() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")

		p := Pokemon{Name: name}
		api.WriteJSON(w, req, p, http.StatusOK)
	}
}
