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
	maxNumberPage = 6
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
		person, err := shikimori_api.GetPersonInfo(i)
		if err != nil {
			log.Printf("Error in person %d: %q", i, err.Error())
		} else {
			log.Printf("Person %d, name %q exists", i, person.Name)
		}
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

	// Смотрим содержимое таблицы студий
	//type RowStudio struct {
	//	id           uuid.UUID
	//	shikimori_id string
	//	studio_name  string
	//	image        string
	//}
	//
	//rows, err = postgres.Conn.Query(context.Background(), "SELECT * FROM public.studio_table")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//
	//var rowSlice1 []RowStudio
	//for rows.Next() {
	//	var r RowStudio
	//	err := rows.Scan(&r.id, &r.shikimori_id, &r.studio_name, &r.image)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	rowSlice1 = append(rowSlice1, r)
	//}
	//if err := rows.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("STUDIO_TABLE$")
	//fmt.Println(rowSlice1)

	//конец просмотра содержимого таблицы
}
