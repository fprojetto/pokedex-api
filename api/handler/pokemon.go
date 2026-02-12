package handler

import (
	"github.com/fprojetto/pokedex-api/pkg/server"

	"net/http"
)

type pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"isLegendary"`
}

func GetPokemon() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")

		p := pokemon{Name: name}
		server.WriteJSON(w, req, p, http.StatusOK)
	}
}

func GetPokemonTranslated() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		name := req.PathValue("name")

		p := pokemon{Name: name}
		server.WriteJSON(w, req, p, http.StatusOK)
	}
}
