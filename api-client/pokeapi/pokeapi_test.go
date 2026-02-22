package pokeapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fprojetto/pokedex-api/api-client/pokeapi"
	"github.com/fprojetto/pokedex-api/model"
	"github.com/fprojetto/pokedex-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPokemonInfo(t *testing.T) {
	tests := []struct {
		name           string
		pokemonName    string
		mockResponse   any
		mockStatus     int
		expectedResult model.Pokemon
		expectedError  error
	}{
		{
			name:        "successful retrieval",
			pokemonName: "pikachu",
			mockStatus:  http.StatusOK,
			mockResponse: pokeapi.PokemonSpeciesResponse{
				ID:          25,
				Name:        "pikachu",
				Habitat:     "forest",
				IsLegendary: pokeapi.BoolPtr(false),
				FlavorTextEntries: []pokeapi.FlavorTextEntry{
					{
						FlavorText: "When several of these POKéMON gather, their electricity could build and cause lightning storms.",
						Language: struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						}{Name: "en"},
					},
				},
			},
			expectedResult: model.Pokemon{
				Name:        "pikachu",
				Description: "When several of these POKéMON gather, their electricity could build and cause lightning storms.",
				Habitat:     "forest",
				IsLegendary: pokeapi.BoolPtr(false),
			},
			expectedError: nil,
		},
		{
			name:        "no legendary data",
			pokemonName: "pikachu",
			mockStatus:  http.StatusOK,
			mockResponse: pokeapi.PokemonSpeciesResponse{
				ID:      25,
				Name:    "pikachu",
				Habitat: "forest",
				FlavorTextEntries: []pokeapi.FlavorTextEntry{
					{
						FlavorText: "When several of these POKéMON gather, their electricity could build and cause lightning storms.",
						Language: struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						}{Name: "en"},
					},
				},
			},
			expectedResult: model.Pokemon{},
			expectedError:  service.ErrMissingData,
		},
		{
			name:        "no english description",
			pokemonName: "pikachu",
			mockStatus:  http.StatusOK,
			mockResponse: pokeapi.PokemonSpeciesResponse{
				Name: "pikachu",
				FlavorTextEntries: []pokeapi.FlavorTextEntry{
					{
						FlavorText: "Pikachu test description",
						Language: struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						}{Name: "fr"},
					},
				},
			},
			expectedResult: model.Pokemon{},
			expectedError:  service.ErrMissingData,
		},
		{
			name:        "api error",
			pokemonName: "pikachu",
			mockStatus:  http.StatusInternalServerError,
			mockResponse: map[string]string{
				"error": "internal server error",
			},
			expectedError: service.ErrServiceUnavailable,
		},
		{
			name:        "api not found",
			pokemonName: "pikachu",
			mockStatus:  http.StatusNotFound,
			mockResponse: map[string]string{
				"error": "not found",
			},
			expectedError: service.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/api/v2/pokemon-species/" + tt.pokemonName
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				w.WriteHeader(tt.mockStatus)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer ts.Close()

			client, err := pokeapi.NewClient(ts.URL)
			if err != nil {
				require.NoError(t, err, "Failed to create pokeapi client")
			}
			result, err := client.PokemonInfo(context.Background(), tt.pokemonName)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.Habitat, result.Habitat)
				assert.Equal(t, tt.expectedResult.IsLegendary, result.IsLegendary)
			}
		})
	}
}
