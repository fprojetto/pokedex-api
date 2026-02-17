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

	"github.com/fprojetto/pokedex-api/application"
	"github.com/fprojetto/pokedex-api/config"
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

func TestGetPokemonE2EInProcess(t *testing.T) {
	// 1. Setup Mock PokeAPI
	mockPokeAPI := runMockPokeAPI()
	defer mockPokeAPI.Close()

	// 2. Setup environment vars for application
	os.Setenv("POKEMON_API_URL", mockPokeAPI.URL)
	appPort := "9090"
	os.Setenv("PORT", appPort) // Random available port
	defer os.Unsetenv("POKEMON_API_URL")
	defer os.Unsetenv("PORT")

	cfg, err := config.New()
	require.NoError(t, err)

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
		t.Fatal("Application failed to start")
	}

	// 3. Perform Request
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
	assert.NotEmpty(t, result.Data.IsLegendary)

	ctx.Done()
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
