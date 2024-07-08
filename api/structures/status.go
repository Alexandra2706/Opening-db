package structures

import (
	"encoding/json"
	"errors"
)

type AnimeStatus string

const (
	AnimeStatusANONS    AnimeStatus = "anons"
	AnimeStatusONGOING  AnimeStatus = "ongoing"
	AnimeStatusRELEASED AnimeStatus = "released"
)

var animeStatusValidator = map[AnimeStatus]bool{
	AnimeStatusANONS:    true,
	AnimeStatusONGOING:  true,
	AnimeStatusRELEASED: true,
}

func (a AnimeStatus) MarshalJSON() ([]byte, error) {
	// Валидация
	_, ok := animeStatusValidator[a]
	if !ok {
		return []byte{}, errors.New("Invalid anime status")
	}

	return json.Marshal(string(a))
}

func (a *AnimeStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	// Валидация
	_, ok := animeStatusValidator[AnimeStatus(s)]
	if !ok {
		return errors.New("Invalid anime status")
	}

	*a = AnimeStatus(s)
	return nil
}
