package main

import (
	"log"
	shikimori_api "shikimori-sync/shikimori-api"
)

const (
	maxNumberPage = 5
)

func main() {

	for i := 1; i < maxNumberPage; i++ {
		anime, err := shikimori_api.GetAnimeInfo(i)
		if err != nil {
			log.Printf("Error in anime %d: %q", i, err.Error())
		} else {
			log.Printf("Anime %d, called %q exists", i, anime.Name)
		}
	}

	anime, err := shikimori_api.ListAnime()
	if err != nil {
		log.Printf("Error in list anime: %q", err.Error())
	} else {
		log.Println("List of anime ids:", anime)
	}

}
