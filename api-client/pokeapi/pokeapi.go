package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fprojetto/pokedex-api/model"
	"github.com/fprojetto/pokedex-api/service"
)

type PokemonSpeciesResponse struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Habitat           string            `json:"habitat"`
	IsLegendary       *bool             `json:"is_legendary"`
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
}

type FlavorTextEntry struct {
	FlavorText string `json:"flavor_text"`
	Language   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"language"`
}

type PokemonClient struct {
	pokeAPIURL string
	client     *http.Client
}

func NewClient(pokeAPIURL string) (*PokemonClient, error) {
	if pokeAPIURL == "" {
		return nil, errors.New("pokeAPIURL empty string")
	}

	return &PokemonClient{
		pokeAPIURL: pokeAPIURL,
		client:     httpClient(),
	}, nil
}

func (c *PokemonClient) PokemonInfo(ctx context.Context, name string) (model.Pokemon, error) {
	res, err := c.getBasicInfo(ctx, name)
	if err != nil {
		return model.Pokemon{}, errors.Join(err, service.ErrServiceUnavailable)
	}

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return model.Pokemon{}, service.ErrNotFound
		default:
			return model.Pokemon{}, service.ErrServiceUnavailable

		}
	}

	var species PokemonSpeciesResponse
	if err := json.NewDecoder(res.Body).Decode(&species); err != nil {
		return model.Pokemon{}, errors.Join(err, service.ErrServiceUnavailable)
	}

	if species.IsLegendary == nil {
		return model.Pokemon{}, service.ErrMissingData
	}

	engDesc := tryToFindEnglishDescription(species.FlavorTextEntries)

	return model.Pokemon{
		Name:        species.Name,
		Description: engDesc,
		Habitat:     species.Habitat,
		IsLegendary: species.IsLegendary,
	}, nil
}

func tryToFindEnglishDescription(entries []FlavorTextEntry) string {
	for e := range entries {
		if strings.ToLower(entries[e].Language.Name) == "en" {
			return entries[e].FlavorText
		}
	}

	return ""
}

func (c *PokemonClient) getBasicInfo(ctx context.Context, name string) (*http.Response, error) {
	getPokemonURL := fmt.Sprintf("%s/api/v2/pokemon-species/%s", c.pokeAPIURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", getPokemonURL, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func BoolPtr(b bool) *bool {
	return &b
}

func httpClient() *http.Client {
	// Define the Transport (Network Layer)
	t := &http.Transport{
		// 1. Connection Dialing settings
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // Max time to establish a TCP connection
			KeepAlive: 30 * time.Second, // Probe interval for active connections
		}).DialContext,

		// 2. TLS/SSL Handshake
		TLSHandshakeTimeout: 10 * time.Second, // Max time for HTTPS handshake

		// 3. Connection Pooling (Very Important!)
		MaxIdleConns:        100,              // Total max idle connections across all hosts
		MaxIdleConnsPerHost: 100,              // Max idle connections for a SINGLE host (Default is 2!)
		IdleConnTimeout:     90 * time.Second, // How long an idle connection is kept open

		// 4. Response Reading
		ResponseHeaderTimeout: 10 * time.Second, // Max time to wait for server response headers
	}

	// Define the Client (High Level)
	client := &http.Client{
		Transport: t,

		// 5. Total Request Timeout
		// This includes Dial, TLS, sending request, waiting for response, and reading body.
		Timeout: 30 * time.Second,
	}

	return client
}
