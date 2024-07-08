package structures

import (
	"encoding/json"
	"errors"
)

type AnimeKind string

const (
	AnimeKindTV    AnimeKind = "tv"
	AnimeKindMOVIE AnimeKind = "movie"
	AnimeKindOVA   AnimeKind = "ova"
	AnimeKindONA   AnimeKind = "ona"
	AnimeKindOTHER AnimeKind = "other"
)

var animeKindValidator = map[AnimeKind]bool{
	AnimeKindTV:    true,
	AnimeKindMOVIE: true,
	AnimeKindOVA:   true,
	AnimeKindONA:   true,
	AnimeKindOTHER: true,
}

func (a AnimeKind) MarshalJSON() ([]byte, error) {
	// Валидация
	_, ok := animeKindValidator[a]
	if !ok {
		return []byte{}, errors.New("Invalid anime kind")
	}

	return json.Marshal(string(a))
}

func (a *AnimeKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	// Валидация
	_, ok := animeKindValidator[AnimeKind(s)]
	if !ok {
		return errors.New("Invalid anime kind")
	}

	*a = AnimeKind(s)
	return nil
}
