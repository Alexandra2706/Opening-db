package postgres

import (
	"api/structures"
	"context"
)

func ListAnime(limit int, offset int) ([]structures.AnimeInfo, error) {

	rows, err := Conn.Query(context.Background(), `
		SELECT id, anime_name, name_russian, name_english, 
		       name_japanese, name_synonyms, anime_status, episodes, 
		       episodes_aired, aired_on, released_on, duration, 
		       updated_at, next_episode_at, 
		       image, genres, studios, videos, 
		       screenshots, shikimori_id, shikimori_kind, shikimori_rating, 
		       shikimori_description, shikimori_description_html, shikimori_last_revision, myanimelist_id,
		       myanimelist_score
		FROM public.animes LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowSlice []structures.AnimeInfo
	for rows.Next() {
		var r structures.AnimeInfo
		err := rows.Scan(
			&r.Id, &r.AnimeName, &r.NameRussian, &r.NameEnglish,
			&r.NameJapanese, &r.NameSynonyms, &r.AnimeStatus.Value, &r.Episodes,
			&r.EpisodesAired, &r.AiredOn.T, &r.ReleasedOn.T, &r.Duration,
			&r.UpdatedAt.T, &r.NextEpisodeAt,
			&r.Image, &r.Genres, &r.Studios, &r.Videos,
			&r.Screenshots, &r.ShikimoriId, &r.ShikimoriKind.Value, &r.ShikimoriRating.Value,
			&r.ShikimoriDescription, &r.ShikimoriDescriptionHtml, &r.ShikimoriLastRevision.T, &r.MyanimelistId,
			&r.MyanimelistScore)
		if err != nil {
			return nil, err
		}
		rowSlice = append(rowSlice, r)
	}

	return rowSlice, nil
}
