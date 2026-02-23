package service

import (
	"context"
	"errors"

	"github.com/fprojetto/pokedex-api/internal/model"
)

var (
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNotFound           = errors.New("pokemon not found")
	ErrMissingData        = errors.New("pokemon data is missing")
)

type TranslationStyle string

const (
	Yoda        TranslationStyle = "yoda"
	Shakespeare                  = "shakespeare"
)

type PokemonInfoGetter func(ctx context.Context, name string) (model.Pokemon, error)
type Translator func(ctx context.Context, translationStyle TranslationStyle, text string) (string, error)

func PokemonGetterService(getter PokemonInfoGetter) func(ctx context.Context, name string) (model.Pokemon, error) {
	return func(ctx context.Context, name string) (model.Pokemon, error) {
		p, err := getter(ctx, name)
		if err != nil {
			return model.Pokemon{}, err
		}

		if err := validate(p); err != nil {
			return model.Pokemon{}, err
		}

		return p, nil
	}
}

func PokemonGetterTranslatorService(
	getter PokemonInfoGetter,
	translator Translator,
) func(ctx context.Context, name string) (model.Pokemon, error) {
	getterService := PokemonGetterService(getter)
	translatorService := pokemonTranslatorService(translator)
	return func(ctx context.Context, name string) (model.Pokemon, error) {
		p, err := getterService(ctx, name)
		if err != nil {
			return model.Pokemon{}, err
		}

		return translatorService(ctx, p), nil
	}
}

func pokemonTranslatorService(translator Translator) func(ctx context.Context, p model.Pokemon) model.Pokemon {
	return func(ctx context.Context, p model.Pokemon) model.Pokemon {
		var translationStyle TranslationStyle
		if p.Habitat == "cave" || (p.IsLegendary != nil && *p.IsLegendary) {
			translationStyle = Yoda
		} else {
			translationStyle = Shakespeare
		}

		if translatedDescription, err := translator(ctx, translationStyle, p.Description); err == nil {
			p.Description = translatedDescription
		}

		return p
	}
}

func validate(p model.Pokemon) error {
	if p.Name == "" || p.Description == "" || p.Habitat == "" || p.IsLegendary == nil {
		return ErrMissingData
	}
	return nil
}
