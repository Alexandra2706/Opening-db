package structures

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

type AnimeStatus string
type AnimeKind string
type AnimeRating string

const (
	AnimeStatusANONS    AnimeStatus = "anons"
	AnimeStatusONGOING  AnimeStatus = "ongoing"
	AnimeStatusRELEASED AnimeStatus = "released"
)

const (
	AnimeKindTV    AnimeKind = "tv"
	AnimeKindMOVIE AnimeKind = "movie"
	AnimeKindOVA   AnimeKind = "ova"
	AnimeKindONA   AnimeKind = "ona"
	AnimeKindOTHER AnimeKind = "other"
)

const (
	AnimeRatingRPLUS AnimeRating = "r_plus"
	AnimeRatingPG13  AnimeRating = "pg_13"
	AnimeRatingR     AnimeRating = "r"
	AnimeRatingG     AnimeRating = "g"
	AnimeRatingRX    AnimeRating = "rx"
	AnimeRatingPG    AnimeRating = "pg"
)

var enumValidator = map[string]map[any]bool{
	"AnimeStatus": {
		AnimeStatusANONS:    true,
		AnimeStatusONGOING:  true,
		AnimeStatusRELEASED: true,
	},
	"AnimeKind": {
		AnimeKindTV:    true,
		AnimeKindMOVIE: true,
		AnimeKindOVA:   true,
		AnimeKindONA:   true,
		AnimeKindOTHER: true,
	},
	"AnimeRating": {
		AnimeRatingRPLUS: true,
		AnimeRatingPG13:  true,
		AnimeRatingR:     true,
		AnimeRatingG:     true,
		AnimeRatingRX:    true,
		AnimeRatingPG:    true,
	},
}

type ValidatableEnum[T comparable] struct {
	Value T
}

func (v *ValidatableEnum[T]) Validate(V map[any]bool) error {
	_, ok := V[v.Value]
	log.Println("Validate v.Value=", v.Value)
	if !ok {
		return errors.New("Invalid validation")
	}
	return nil
}

func (v *ValidatableEnum[T]) MarshalJSON() ([]byte, error) {
	// Валидация
	if v.Validate(enumValidator[reflect.TypeOf(v.Value).Name()]) != nil {
		return []byte{}, errors.New("Invalid validation")
	}
	return json.Marshal(v.Value)
}

func (v *ValidatableEnum[T]) UnmarshalJSON(data []byte) error {
	var s T
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	v.Value = s
	if v.Validate(enumValidator[reflect.TypeOf(s).Name()]) != nil {
		return errors.New("Invalid validation")
	}

	return nil
}
