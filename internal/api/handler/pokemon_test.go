package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fprojetto/pokedex-api/internal/api"
	"github.com/fprojetto/pokedex-api/internal/api/handler"
	"github.com/fprojetto/pokedex-api/internal/model"
	"github.com/fprojetto/pokedex-api/internal/service"
	"github.com/fprojetto/pokedex-api/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type pokemonServiceMock struct {
	mock.Mock
}

func (m *pokemonServiceMock) GetPokemon(ctx context.Context, name string) (model.Pokemon, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(model.Pokemon), args.Error(1)
}

func (m *pokemonServiceMock) GetPokemonTranslated(ctx context.Context, name string) (model.Pokemon, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(model.Pokemon), args.Error(1)
}

func TestGetPokemon(t *testing.T) {
	tests := []struct {
		name                  string
		pokemonName           string
		mockReturnPokemon     model.Pokemon
		mockReturnError       error
		expectedStatusCode    int
		expectedBodyAssertion func(t *testing.T, body []byte, expectedPokemon model.Pokemon)
	}{
		{
			name:        "GET /api/pokemon/{name} success",
			pokemonName: "mewtwo",
			mockReturnPokemon: model.Pokemon{
				Name:        "mewtwo",
				Description: "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
				Habitat:     "rare",
				IsLegendary: client.BoolPtr(false),
			},
			mockReturnError:    nil,
			expectedStatusCode: http.StatusOK,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope struct {
					Data handler.Pokemon `json:"data"`
				}
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")

				assert.Equal(t, expectedPokemon.Name, envelope.Data.Name)
				assert.Equal(t, expectedPokemon.Description, envelope.Data.Description)
				assert.Equal(t, expectedPokemon.Habitat, envelope.Data.Habitat)
				assert.Equal(t, expectedPokemon.IsLegendary, envelope.Data.IsLegendary)
			},
		},
		{
			name:               "GET /api/pokemon/{name} internal error",
			pokemonName:        "bulbasaur",
			mockReturnPokemon:  model.Pokemon{},
			mockReturnError:    service.ErrServiceUnavailable,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope api.Envelope
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")
				assert.NotNil(t, envelope.Error)
			},
		},
		{
			name:               "GET /api/pokemon/{name} not found",
			pokemonName:        "pikachu",
			mockReturnPokemon:  model.Pokemon{},
			mockReturnError:    service.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope api.Envelope
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")
				assert.NotNil(t, envelope.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &pokemonServiceMock{}
			getPokemonHandler := handler.GetPokemon(mockService.GetPokemon)

			mockService.On("GetPokemon", mock.Anything, tt.pokemonName).Return(tt.mockReturnPokemon, tt.mockReturnError)

			req := httptest.NewRequest("GET", "/api/pokemon/"+tt.pokemonName, nil)
			req.SetPathValue("name", tt.pokemonName)
			res := httptest.NewRecorder()
			getPokemonHandler(res, req)

			assert.Equal(t, tt.expectedStatusCode, res.Code)

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err, "failed to read response body")

			tt.expectedBodyAssertion(t, body, tt.mockReturnPokemon)
		})
	}
}

func TestGetPokemonTranslated(t *testing.T) {
	tests := []struct {
		name                  string
		pokemonName           string
		mockReturnPokemon     model.Pokemon
		mockReturnError       error
		expectedStatusCode    int
		expectedBodyAssertion func(t *testing.T, body []byte, expectedPokemon model.Pokemon)
	}{
		{
			name:        "GET /api/pokemon/translated/{name} success",
			pokemonName: "mewtwo",
			mockReturnPokemon: model.Pokemon{
				Name:        "mewtwo",
				Description: "Translated description.",
				Habitat:     "rare",
				IsLegendary: client.BoolPtr(true),
			},
			mockReturnError:    nil,
			expectedStatusCode: http.StatusOK,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope struct {
					Data handler.Pokemon `json:"data"`
				}
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")

				assert.Equal(t, expectedPokemon.Name, envelope.Data.Name)
				assert.Equal(t, expectedPokemon.Description, envelope.Data.Description)
				assert.Equal(t, expectedPokemon.Habitat, envelope.Data.Habitat)
				assert.Equal(t, expectedPokemon.IsLegendary, envelope.Data.IsLegendary)
			},
		},
		{
			name:               "GET /api/pokemon/translated/{name} internal error",
			pokemonName:        "bulbasaur",
			mockReturnPokemon:  model.Pokemon{},
			mockReturnError:    service.ErrServiceUnavailable,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope api.Envelope
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")
				assert.NotNil(t, envelope.Error)
			},
		},
		{
			name:               "GET /api/pokemon/translated/{name} not found",
			pokemonName:        "pikachu",
			mockReturnPokemon:  model.Pokemon{},
			mockReturnError:    service.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedBodyAssertion: func(t *testing.T, body []byte, expectedPokemon model.Pokemon) {
				var envelope api.Envelope
				err := json.Unmarshal(body, &envelope)
				require.NoError(t, err, "failed to unmarshal response")
				assert.NotNil(t, envelope.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &pokemonServiceMock{}
			getPokemonTranslatedHandler := handler.GetPokemonTranslated(mockService.GetPokemonTranslated)

			mockService.On("GetPokemonTranslated", mock.Anything, tt.pokemonName).Return(tt.mockReturnPokemon, tt.mockReturnError)

			req := httptest.NewRequest("GET", "/api/pokemon/translated/"+tt.pokemonName, nil)
			req.SetPathValue("name", tt.pokemonName)
			res := httptest.NewRecorder()
			getPokemonTranslatedHandler(res, req)

			assert.Equal(t, tt.expectedStatusCode, res.Code)

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err, "failed to read response body")

			tt.expectedBodyAssertion(t, body, tt.mockReturnPokemon)
		})
	}
}
