package main

import (
	"log"
	"time"

	shikimori_api "shikimori-sync/shikimori-api"
)

const (
	maxNumberPage = 10
)

func main() {

	for i := 1; i < maxNumberPage; i++ {
		anime, err := shikimori_api.GetAnimeInfo(i)
		if err != nil {
			log.Printf("Error in anime %d: %q", i, err.Error())
		} else {
			log.Printf("Anime %d, called %q exists", i, anime.Name)
		}
		time.Sleep(2 * time.Second)
	}

}
