package service

import (
	"context"
	"errors"

	"github.com/fprojetto/pokedex-api/model"
)

var (
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNotFound           = errors.New("pokemon not found")
	ErrMissingData        = errors.New("pokemon data is missing")
)

type PokemonInfoGetter func(ctx context.Context, name string) (model.Pokemon, error)

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

func validate(p model.Pokemon) error {
	if p.Name == "" || p.Description == "" || p.Habitat == "" || p.IsLegendary == nil {
		return ErrMissingData
	}
	return nil
}
