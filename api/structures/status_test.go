package structures

import (
	"log"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	data := []byte("\"released\"")
	log.Println("data=", data)
	log.Println("data=", string(data))

	v := &ValidatableEnum[AnimeStatus]{}
	log.Println("v0=", v)
	err := v.UnmarshalJSON(data)
	log.Println("v1=", v)
	if err != nil {
		t.Errorf("Should not produce an error, but have one: %q", err.Error())
	}
}
