package shikimori_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	baseEndpoint = "https://shikimori.one/api/"
	debug        = true
)

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

	time.Sleep(1000 * time.Millisecond)
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

func ListAnime() ([]int, error) {
	ids := make([]int, 0)
	i := 0
	for true {
		animes := &[]AnimeShort{}
		err := MakeRequest("animes/?limit=50&page="+strconv.Itoa(i)+"&order=id", animes)
		if err != nil {
			return nil, err
		}
		if len(*animes) == 0 || (debug && i > 0) {
			return ids, nil
		}
		for _, short := range *animes {
			ids = append(ids, short.Id)
		}
		i++
	}
	return nil, nil
}
