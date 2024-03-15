package shikimori_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const baseEndpoint = "https://shikimori.one/api/"

func MakeRequest(currentUrl string, result interface{}) error {
	resp, err := http.Get(baseEndpoint + currentUrl)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot read body: %q", err.Error()))
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func GetAnimeInfo(id int) (*Anime, error) {
	anime := &Anime{}
	err := MakeRequest("animes/"+strconv.Itoa(id), anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}
