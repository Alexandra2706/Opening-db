package main

import (
	"context"
	"fmt"
	"log"
	"shikimori-sync/logic"
	"shikimori-sync/postgres"
	shikimori_api "shikimori-sync/shikimori-api"
)

const (
	maxNumberPage = 30
	imageEndpoint = "https://shikimori.one"
)

var kindsMap map[string]string = map[string]string{
	"tv":    "tv",
	"movie": "movie",
	"ova":   "ova",
	"ona":   "ona",
	"other": "other",
}

var ratingMap map[string]string = map[string]string{
	"r_plus": "r_plus",
	"pg_13":  "pg_13",
	"r":      "r",
	"g":      "g",
	"rx":     "rx",
	"pg":     "pg",
}

func main() {
	defer postgres.CloseConnection()

	animeList, err := shikimori_api.ListAnime()
	if err != nil {
		log.Fatalf("Error in list anime: %q", err.Error())
	} else {
		log.Println("List of anime ids:", animeList)
	}

	//временно: удаляем все записи из таблицы жанров
	//_, err = postgres.Conn.Exec(context.Background(), "delete from genres_table")
	//if err != nil {
	//	log.Println(err)
	//}

	for _, aid := range animeList {

		logic.CreateOrUpdateAnime(aid)

	}

	for i := 1; i < maxNumberPage; i++ {
		logic.CreateOrUpdatePerson(i)
	}

	// Смотрим содержимое таблицы жанров
	type Row struct {
		id           string
		shikimori_id string
		genre_name   string
		russian      string
	}

	rows, err := postgres.Conn.Query(context.Background(), "SELECT * FROM genres_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rowSlice []Row
	for rows.Next() {
		var r Row
		err := rows.Scan(&r.id, &r.shikimori_id, &r.genre_name, &r.russian)
		if err != nil {
			log.Fatal(err)
		}
		rowSlice = append(rowSlice, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("GENRE_TABLE$")
	fmt.Println(rowSlice)

	//конец просмотра содержимого таблицы
}
