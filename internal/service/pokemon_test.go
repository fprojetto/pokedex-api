package service

import (
	"context"
	"errors"
	"testing"

	"github.com/fprojetto/pokedex-api/internal/model"
	"github.com/fprojetto/pokedex-api/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestPokemonGetterService(t *testing.T) {
	testCases := []struct {
		name          string
		mockGetter    PokemonInfoGetter
		pokemonName   string
		expected      model.Pokemon
		expectedError error
	}{
		{
			name: "Success - Valid Pokemon",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) {
				isLegendary := false
				return model.Pokemon{
					Name:        "charmander",
					Description: "A small orange lizard pokemon.",
					Habitat:     "mountain",
					IsLegendary: &isLegendary,
				}, nil
			},
			pokemonName: "charmander",
			expected: model.Pokemon{
				Name:        "charmander",
				Description: "A small orange lizard pokemon.",
				Habitat:     "mountain",
				IsLegendary: client.BoolPtr(false),
			},
			expectedError: nil,
		},
		{
			name: "Error - Pokemon Not Found",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				return model.Pokemon{}, ErrNotFound
			},
			pokemonName:   "nonexistent",
			expected:      model.Pokemon{},
			expectedError: ErrNotFound,
		},
		{
			name: "Error - Missing Data (empty name)",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := false
				return model.Pokemon{
					Name:        "",
					Description: "A small orange lizard pokemon.",
					Habitat:     "mountain",
					IsLegendary: &isLegendary,
				}, nil
			},
			pokemonName:   "invalid",
			expected:      model.Pokemon{},
			expectedError: ErrMissingData,
		},
		{
			name: "Error - Missing Data (nil isLegendary)",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				return model.Pokemon{
					Name:        "charmander",
					Description: "A small orange lizard pokemon.",
					Habitat:     "mountain",
					IsLegendary: nil,
				}, nil
			},
			pokemonName:   "invalid",
			expected:      model.Pokemon{},
			expectedError: ErrMissingData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			service := PokemonGetterService(tc.mockGetter) // Use tc.mockGetter directly
			result, err := service(ctx, tc.pokemonName)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, tc.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.Description, result.Description)
				assert.Equal(t, tc.expected.Habitat, result.Habitat)
				assert.Equal(t, *tc.expected.IsLegendary, *result.IsLegendary)
			}
		})
	}
}

func TestPokemonGetterTranslatorService(t *testing.T) {
	testCases := []struct {
		name           string
		mockGetter     PokemonInfoGetter // Use the function type directly
		mockTranslator Translator        // Use the function type directly
		pokemonName    string
		expected       model.Pokemon
		expectedError  error
	}{
		{
			name: "Success - Yoda Translation for Legendary Pokemon",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := true
				return model.Pokemon{
					Name:        "mewtwo",
					Description: "A legendary psychic pokemon.",
					Habitat:     "rare",
					IsLegendary: &isLegendary,
				}, nil
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				assert.Equal(t, string(Yoda), string(style))
				return "A legendary psychic pokemon, hmmm.", nil
			},
			pokemonName: "mewtwo",
			expected: model.Pokemon{
				Name:        "mewtwo",
				Description: "A legendary psychic pokemon, hmmm.",
				Habitat:     "rare",
				IsLegendary: client.BoolPtr(true),
			},
			expectedError: nil,
		},
		{
			name: "Success - Yoda Translation for Cave Habitat Pokemon",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := false
				return model.Pokemon{
					Name:        "zubat",
					Description: "A bat pokemon that lives in caves.",
					Habitat:     "cave",
					IsLegendary: &isLegendary,
				}, nil
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				assert.Equal(t, string(Yoda), string(style))
				return "A bat pokemon that lives in caves, hmmm.", nil
			},
			pokemonName: "zubat",
			expected: model.Pokemon{
				Name:        "zubat",
				Description: "A bat pokemon that lives in caves, hmmm.",
				Habitat:     "cave",
				IsLegendary: client.BoolPtr(false),
			},
			expectedError: nil,
		},
		{
			name: "Success - Shakespeare Translation for Other Pokemon",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := false
				return model.Pokemon{
					Name:        "pikachu",
					Description: "A small electric mouse pokemon.",
					Habitat:     "forest",
					IsLegendary: &isLegendary,
				}, nil
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				assert.Equal(t, string(Shakespeare), string(style))
				return "A small electric mouse pokemon, forsooth.", nil
			},
			pokemonName: "pikachu",
			expected: model.Pokemon{
				Name:        "pikachu",
				Description: "A small electric mouse pokemon, forsooth.",
				Habitat:     "forest",
				IsLegendary: client.BoolPtr(false),
			},
			expectedError: nil,
		},
		{
			name: "Success - Translation Fails, Original Description Returned",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := false
				return model.Pokemon{
					Name:        "squirtle",
					Description: "A small turtle pokemon.",
					Habitat:     "waters-edge",
					IsLegendary: &isLegendary,
				}, nil
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				return "", errors.New("translation service unavailable")
			},
			pokemonName: "squirtle",
			expected: model.Pokemon{
				Name:        "squirtle",
				Description: "A small turtle pokemon.", // Original description expected
				Habitat:     "waters-edge",
				IsLegendary: client.BoolPtr(false),
			},
			expectedError: nil,
		},
		{
			name: "Error - Pokemon Not Found by Getter",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				return model.Pokemon{}, ErrNotFound
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				return "", nil // Should not be called
			},
			pokemonName:   "nonexistent",
			expected:      model.Pokemon{},
			expectedError: ErrNotFound,
		},
		{
			name: "Error - Missing Data from Getter",
			mockGetter: func(ctx context.Context, name string) (model.Pokemon, error) { // Assign the function directly
				isLegendary := false
				return model.Pokemon{
					Name:        "", // Missing name
					Description: "A small electric mouse pokemon.",
					Habitat:     "forest",
					IsLegendary: &isLegendary,
				}, nil
			},
			mockTranslator: func(ctx context.Context, style TranslationStyle, text string) (string, error) { // Assign the function directly
				return "", nil // Should not be called
			},
			pokemonName:   "invalid",
			expected:      model.Pokemon{},
			expectedError: ErrMissingData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			service := PokemonGetterTranslatorService(tc.mockGetter, tc.mockTranslator) // Use tc.mockGetter and tc.mockTranslator directly
			result, err := service(ctx, tc.pokemonName)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, tc.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.Description, result.Description)
				assert.Equal(t, tc.expected.Habitat, result.Habitat)
				assert.Equal(t, *tc.expected.IsLegendary, *result.IsLegendary)
			}
		})
	}
}
