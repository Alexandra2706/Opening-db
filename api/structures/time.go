package structures

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type Time struct {
	T *time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.T == nil {
		return []byte("null"), nil
	}
	tm := t.T.UTC().Round(time.Second)
	tmJSON, err := json.Marshal(tm)
	if err != nil {
		return []byte{}, errors.New("Invalid anime marshall time")
	}
	log.Println(string(tmJSON))
	return tmJSON, nil
}
