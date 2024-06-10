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
		/*
			anime, err := shikimori_api.GetAnimeInfo(aid)
			if err != nil {
				log.Printf("Error in anime %d: %q", aid, err.Error())
				continue
			}
			log.Printf("Saving anime %d, called %q", aid, anime.Name)

			// Добавляем image в s3 и images_table, если все ок, то получаем url для таблица анимэ, иначе записываем nil в таблицу анимэ
			err = s3.CreateOrUpdateImage(imageEndpoint + anime.Image.Original)
			if err != nil {
				log.Fatalf("Error in anime image: %q", err.Error())
			}

			var genresList = []string{}
			for i := 0; i < len(anime.Genres); i++ {
				res := postgres.CreateOrUpdateGenre(anime.Genres[i].Id, anime.Genres[i].Name, anime.Genres[i].Russian)
				genresList = append(genresList, res)
			}

			var studioList = []uuid.UUID{}
			for i := 0; i < len(anime.Studios); i++ {
				res := logic.CreateOrUpdateStudio(anime.Studios[i].Id, anime.Studios[i].Name, imageEndpoint+anime.Studios[i].Image)
				studioList = append(studioList, res)
			}
			fmt.Println(studioList)

			// видео заполнить потом
			var videoList = []string{}

			var screenshotsList = []string{}
			for i := 0; i < len(anime.Screenshots); i++ {
				err := s3.CreateOrUpdateImage(imageEndpoint + anime.Screenshots[i].Original)
				if err != nil {
					log.Fatalf("Error in screenshort image: %q", err.Error())
				}
				screenshotsList = append(screenshotsList, imageEndpoint+anime.Screenshots[i].Original)
			}

			//выбор kind
			dbKind, ok := kindsMap[anime.Kind]
			if !ok {
				log.Printf("Anime kind not found: %q", anime.Kind)
				dbKind = "other"
			}

			dbRating, ok := ratingMap[anime.Rating]
			if !ok {
				log.Printf("Anime raiting not found: %q", anime.Rating)
				dbRating = "other" // надо ли тут добавить в enum 'other'?
			}

			row := postgres.Conn.QueryRow(context.Background(), `
				INSERT INTO public.animes (
					anime_name,    name_russian,   name_english,     name_japanese,
				    name_synonyms, anime_status,   episodes,         episodes_aired,
				    aired_on,      released_on,    duration,         licensors_ru,
				    franchise,     updated_at,     next_episode_at,  image,
				    genres,        studios,        videos,           screenshots,
				    shikimori_id,  shikimori_kind, shikimori_rating, shikimori_description,
				    shikimori_description_html, shikimori_last_revision, myanimelist_id, myanimelist_score
				    )
				VALUES (
				    $1, $2, $3, $4,
				    $5, $6, $7, $8,
				    $9, $10, $11, $12,
				    $13, $14, $15, $16,
				    $17, $18, $19, $20,
				    $21, $22::kind, $23::rating, $24,
				    $25, $26, $27, $28) RETURNING id;`,
				anime.Name, anime.Russian, anime.English, anime.Japanese,
				anime.Synonyms, anime.Status, anime.Episodes, anime.EpisodesAired,
				anime.AiredOn, anime.ReleasedOn, anime.Duration, nil,
				nil, anime.UpdatedAt, anime.NextEpisodeAt, imageEndpoint+anime.Image.Original,
				genresList, studioList, videoList, screenshotsList,
				anime.Id, dbKind, dbRating, anime.Description,
				anime.DescriptionHtml, anime.UpdatedAt, anime.MyanimelistId, anime.Score,
			)

			var id string
			err = row.Scan(&id)
			if err != nil {
				log.Printf("Error when saving anime %d: %q", aid, err.Error())
				continue
			}
			log.Printf("Anime %q saved!", anime.Name)
		*/
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

	//конец просмотра содержимого таблицы

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
