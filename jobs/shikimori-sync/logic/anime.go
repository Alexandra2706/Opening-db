package logic

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

func CreateOrUpdateAnime(aid int) {
	anime, err := shikimori_api.GetAnimeInfo(aid)
	if err != nil {
		log.Fatalf("Error in anime %d: %q", aid, err.Error())
		//continue
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
		res := CreateOrUpdateStudio(anime.Studios[i].Id, anime.Studios[i].Name, imageEndpoint+anime.Studios[i].Image)
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
				anime_name,     name_russian,     name_english,          name_japanese, 
			    name_synonyms,  anime_status,     episodes,              episodes_aired, 
			    aired_on,       released_on,      duration,              licensors_ru, 
			    franchise,      next_episode_at,  image, 	             genres,
			    studios,        videos,           screenshots,           shikimori_id,  
			    shikimori_kind, shikimori_rating, shikimori_description, shikimori_description_html,
			    shikimori_last_revision, myanimelist_id, myanimelist_score
			    )
			VALUES (
			    $1, $2, $3, $4, 
			    $5, $6, $7, $8, 
			    $9, $10, $11, $12, 
			    $13, $14, $15, $16, 
			    $17, $18, $19, $20,
			    $21::kind, $22::rating, $23, $24,
			    $25, $26, $27) 
			ON CONFLICT (shikimori_id) 
			DO UPDATE SET 
			     anime_name = $1,     name_russian = $2,      name_english = $3, name_japanese = $4, 
			     name_synonyms = $5,  anime_status = $6,      episodes = $7,     episodes_aired = $8,
			     aired_on = $9,       released_on = $10,      duration = $11,    licensors_ru = $12, 
			     franchise = $13,     next_episode_at = $14,  image = $15, 	     genres = $16,
			     studios = $17,       videos = $18,           screenshots = $19, shikimori_id = $20,  
			     shikimori_kind = $21::kind, shikimori_rating = $22::rating, shikimori_description = $23, shikimori_description_html = $24,
			     shikimori_last_revision = $25, myanimelist_id = $26, myanimelist_score = $27,
			     updated_at = NOW()
			RETURNING id;`,
		anime.Name, anime.Russian, anime.English, anime.Japanese,
		anime.Synonyms, anime.Status, anime.Episodes, anime.EpisodesAired,
		anime.AiredOn, anime.ReleasedOn, anime.Duration, nil,
		nil, anime.NextEpisodeAt, imageEndpoint+anime.Image.Original,
		genresList, studioList, videoList, screenshotsList,
		anime.Id, dbKind, dbRating, anime.Description,
		anime.DescriptionHtml, anime.UpdatedAt, anime.MyanimelistId, anime.Score,
	)

	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Fatalf("Error when saving anime %d: %q", aid, err.Error())
		//continue
	}
	log.Printf("Anime %q saved!", anime.Name)
}
