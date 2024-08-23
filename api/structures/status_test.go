package structures

import (
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	data := []byte("\"released\"")

	v := &ValidatableEnum[AnimeStatus]{}
	err := v.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("Should not produce an error, but have one: %q", err.Error())
	}

	data = []byte("\"new format\"")
	err = v.UnmarshalJSON(data)
	if err == nil {
		t.Errorf("Must be error")
	}
}
