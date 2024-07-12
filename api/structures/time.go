package structures

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println(string(tmJSON))
	//marshal_time := time.Now().UTC()
	return tmJSON, nil
}
