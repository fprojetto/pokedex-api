package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fprojetto/pokedex-api/api"
	"github.com/fprojetto/pokedex-api/api/handler"
)

func TestPokemonRoutes(t *testing.T) {
	mux := api.NewPokemonRouter(
		handler.GetPokemon(),
		handler.GetPokemonTranslated(),
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("GET /api/pokemon/{name}", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/pokemon/pikachu")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var envelope api.Envelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// The current implementation returns the pokemon struct in the Data field
		data, ok := envelope.Data.(map[string]any)
		if !ok {
			t.Fatalf("expected Data to be a map; got %T", envelope.Data)
		}

		if data["name"] != "pikachu" {
			t.Errorf("expected pokemon name 'pikachu'; got '%v'", data["name"])
		}
	})

	t.Run("GET /api/pokemon/translated/{name}", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/pokemon/translated/mewtwo")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var envelope api.Envelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		data, ok := envelope.Data.(map[string]any)
		if !ok {
			t.Fatalf("expected Data to be a map; got %T", envelope.Data)
		}

		if data["name"] != "mewtwo" {
			t.Errorf("expected pokemon name 'mewtwo'; got '%v'", data["name"])
		}
	})
}
