package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"shikimori-sync/postgres"
	"shikimori-sync/s3"
	shikimori_api "shikimori-sync/shikimori-api"
)

const (
	maxNumberPage = 6
	imageEndpoint = "https://shikimori.one"
)

func main() {
	defer postgres.CloseConnection()

	animeList, err := shikimori_api.ListAnime()
	if err != nil {
		log.Printf("Error in list anime: %q", err.Error())
	} else {
		log.Println("List of anime ids:", animeList)
	}

	//временно: удаляем все записи из таблицы жанров
	//_, err = postgres.Conn.Exec(context.Background(), "delete from genres_table")
	//if err != nil {
	//	log.Println(err)
	//}

	for _, aid := range animeList {
		anime, err := shikimori_api.GetAnimeInfo(aid)
		if err != nil {
			log.Printf("Error in anime %d: %q", aid, err.Error())
			continue
		}
		log.Printf("Saving anime %d, called %q", aid, anime.Name)

		var genresList = []string{}
		for i := 0; i < len(anime.Genres); i++ {
			res := postgres.CreateOrUpdateGenre(anime.Genres[i].Id, anime.Genres[i].Name, anime.Genres[i].Russian)
			genresList = append(genresList, res)
		}

		var studioList = []uuid.UUID{}
		for i := 0; i < len(anime.Studios); i++ {
			err = s3.CreateOrUpdateImage(imageEndpoint + anime.Studios[i].Image)
			if err != nil {
				log.Printf("Error in get image: %q", err.Error())
			}
			fmt.Println("данные студии", anime.Studios[i].Id, anime.Studios[i].Name, anime.Studios[i].Image)
			res := postgres.CreateOrUpdateStudio(anime.Studios[i].Id, anime.Studios[i].Name, imageEndpoint+anime.Studios[i].Image)
			studioList = append(studioList, res)
		}
		fmt.Println(studioList)

		//	row := postgres.Conn.QueryRow(context.Background(), `
		//	INSERT INTO public.animes (anime_name, name_russian, name_english, name_japanese, name_synonyms, anime_status,
		//	                          episodes, episodes_aired, aired_on, released_on, duration, licensors_ru, franchise,
		//	                          updated_at, next_episode_at, image, genres, studios)
		//	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id;`,
		//		anime.Name, anime.Russian, anime.English, anime.Japanese, anime.Synonyms, anime.Status,
		//		anime.Episodes, anime.EpisodesAired, anime.AiredOn, anime.ReleasedOn, anime.Duration, nil, nil,
		//		anime.UpdatedAt, anime.NextEpisodeAt, nil, genresList, nil)
		//
		//	var id string
		//	err = row.Scan(&id)
		//	if err != nil {
		//		log.Printf("Error when saving anime %d: %q", aid, err.Error())
		//		continue
		//	}
		//	log.Printf("Anime %q saved!", anime.Name)
		//
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
