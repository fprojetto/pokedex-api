//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fprojetto/pokedex-api/internal/application"
	"github.com/fprojetto/pokedex-api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PokemonResponse struct {
	Data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Habitat     string `json:"habitat"`
		IsLegendary bool   `json:"isLegendary"`
	} `json:"data"`
}

const appPort = "9091"

func TestMain(m *testing.M) {
	// 1. Setup Mock Servers
	mockPokeAPI := runMockPokeAPI()
	defer mockPokeAPI.Close()
	mockFunTranslationsAPI := runMockFunTranslationsAPI()
	defer mockFunTranslationsAPI.Close()

	// 2. Setup environment vars for application
	os.Setenv("POKEMON_API_URL", mockPokeAPI.URL)
	os.Setenv("TRANSLATION_API_URL", mockFunTranslationsAPI.URL)
	os.Setenv("PORT", appPort) // Random available port
	defer os.Unsetenv("POKEMON_API_URL")
	defer os.Unsetenv("TRANSLATION_API_URL")
	defer os.Unsetenv("PORT")

	cfg, err := config.New()
	if err != nil {
		log.Printf("failed to start application: %v\n", err)
		os.Exit(1)
	}

	// Since application.Run is blocking, we run it in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := application.Run(ctx, *cfg); err != nil && err != http.ErrServerClosed {
			log.Printf("App stopped with error: %v\n", err)
		} else {
			log.Println("App stopped")
		}
	}()

	// Wait for app to start
	if appStarted := waitForApp(appPort); !appStarted {
		log.Printf("Application failed to start\n")
		os.Exit(1)
	}

	exitCode := m.Run()
	ctx.Done()
	os.Exit(exitCode)
}

func TestE2EInProcess(t *testing.T) {
	t.Run("Get Pokemon", func(t *testing.T) {
		pokemonName := "mentwo"
		url := fmt.Sprintf("http://localhost:%s/api/pokemon/%s", appPort, pokemonName)

		resp, err := http.Get(url)
		require.NoError(t, err, "Failed to send request to API")
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result PokemonResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "Failed to decode response body")

		assert.Equal(t, pokemonName, result.Data.Name)
		assert.NotEmpty(t, result.Data.Description)
		assert.NotEmpty(t, result.Data.Habitat)
		assert.Equal(t, true, result.Data.IsLegendary)
	})
	t.Run("Get Pokemon Translated", func(t *testing.T) {
		pokemonName := "mentwo"
		url := fmt.Sprintf("http://localhost:%s/api/pokemon/translated/%s", appPort, pokemonName)

		resp, err := http.Get(url)
		require.NoError(t, err, "Failed to send request to API")
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result PokemonResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "Failed to decode response body")

		assert.Equal(t, pokemonName, result.Data.Name)
		assert.Equal(t, "Created by a scientist after years of horrific gene splicing and dna engineering experiments, it was.", result.Data.Description)
		assert.Equal(t, "rare", result.Data.Habitat)
		assert.Equal(t, true, result.Data.IsLegendary)

	})
}

func runMockPokeAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/pokemon-species/mentwo" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":           150,
				"name":         "mentwo",
				"habitat":      "rare",
				"is_legendary": true,
				"flavor_text_entries": []map[string]interface{}{
					{
						"flavor_text": "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
						"language": map[string]interface{}{
							"name": "en",
						},
					},
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func runMockFunTranslationsAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/translate/yodish" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": map[string]interface{}{
					"total": 1,
				},
				"contents": map[string]interface{}{
					"translated":  "Created by a scientist after years of horrific gene splicing and dna engineering experiments, it was.",
					"text":        "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
					"translation": "yoda",
				},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func waitForApp(appPort string) bool {
	appStarted := false
	for i := 0; i < 10 && !appStarted; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", appPort))
		if err == nil && resp.StatusCode == http.StatusOK {
			appStarted = true
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return appStarted
}
