package postgres

import (
	"api/structures"
	"context"
	"fmt"
)

func ListAnime(limit int, offset int) ([]structures.AnimeShortInfo, error) {

	rows, err := Conn.Query(context.Background(), `
		SELECT id, anime_name, name_russian, name_english, name_japanese,
    			name_synonyms, anime_status, episodes, episodes_aired, aired_on, 
    			released_on, duration, 
    			updated_at, next_episode_at, image, genres, studios,
    			videos, screenshots,
    			shikimori_id, shikimori_kind 
		FROM public.animes LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowSlice []structures.AnimeShortInfo
	for rows.Next() {
		var r structures.AnimeShortInfo
		err := rows.Scan(&r.Id, &r.AnimeName, &r.NameRussian, &r.NameEnglish, &r.NameJapanese,
			&r.NameSynonyms, &r.AnimeStatus, &r.Episodes, &r.EpisodesAired, &r.AiredOn.T,
			&r.ReleasedOn.T, &r.Duration,
			&r.UpdatedAt.T, &r.NextEpisodeAt, &r.Image, &r.Genres, &r.Studios,
			&r.Videos, &r.Screenshots,
			&r.ShikimoriId, &r.ShikimoriKind)
		if err != nil {
			return nil, err
		}
		rowSlice = append(rowSlice, r)
	}

	fmt.Println("ANIME LIST")
	fmt.Println(rowSlice)

	return rowSlice, nil

}
