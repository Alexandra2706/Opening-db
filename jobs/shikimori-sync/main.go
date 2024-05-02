package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"shikimori-sync/postgres"
	shikimori_api "shikimori-sync/shikimori-api"
	"strings"
)

const (
	maxNumberPage = 6
)

func CheckRowExists(table, column string, value interface{}) (bool, error) {
	query := fmt.Sprintf("SELECT %s FROM %s where %s = $1 limit 1", column, table, column)
	//fmt.Printf("QUERY= %s\n", query)
	row := postgres.Conn.QueryRow(context.Background(), query, value)
	fmt.Printf("ROW= %s\n", value)
	var tmp interface{}
	err := row.Scan(&tmp)
	// Тут проблема sql.ErrNoRows = 'sql: no rows in result set',
	// А у меня err = 'no rows in result set'
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func createOrUpdateGenre(shikimoriId int, genreName string, russian string) string {
	stringId := ""
	tempStr := strings.ToLower(genreName) // заменить пробелы на нижн подчеркивание
	for _, i := range tempStr {
		if i != ' ' {
			stringId += string(i)
		}
		if i == ' ' {
			stringId += "_"
		}
	}

	res, errCheck := CheckRowExists("genres_table", "id", stringId)
	fmt.Printf("RES, ERR = %q,,, %q\n", res, errCheck)
	if res == false && errCheck == nil {
		fmt.Println("For ADD: ", stringId, shikimoriId, genreName, russian)
		_, err := postgres.Conn.Exec(context.Background(), `
			INSERT INTO public.genres_table (id, shikimori_id, genre_name, russian) VALUES ($1, $2, $3, $4);`,
			stringId, shikimoriId, genreName, russian)
		if err != nil {
			log.Fatalf("Error: %q", err.Error())
		}

	}

	return stringId
}

func main() {
	defer postgres.CloseConnection()

	animeList, err := shikimori_api.ListAnime()
	if err != nil {
		log.Printf("Error in list anime: %q", err.Error())
	} else {
		log.Println("List of anime ids:", animeList)
	}

	//временно: удаляем все записи из таблицы жанров
	_, err = postgres.Conn.Exec(context.Background(), "delete from genres_table")
	if err != nil {
		log.Println(err)
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

	fmt.Println("GANRE_TABLE$")
	fmt.Println(rowSlice)

	//конец просмотра содержимого таблицы

	for _, aid := range animeList {
		anime, err := shikimori_api.GetAnimeInfo(aid)
		if err != nil {
			log.Printf("Error in anime %d: %q", aid, err.Error())
			continue
		}
		log.Printf("Saving anime %d, called %q", aid, anime.Name)

		for i := 0; i < len(anime.Genres); i++ {
			res := createOrUpdateGenre(anime.Genres[i].Id, anime.Genres[i].Name, anime.Genres[i].Russian)
			fmt.Println(res)
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

		fmt.Println("GANRE_TABLE$")
		fmt.Println(rowSlice)

		//конец просмотра содержимого таблицы

		//row := postgres.Conn.QueryRow(context.Background(), `
		//INSERT INTO public.animes (anime_name, name_russian, name_english, name_japanese, name_synonyms, anime_status,
		//                           episodes, episodes_aired, aired_on, released_on, duration, licensors_ru, franchise,
		//                           updated_at, next_episode_at, image, genres)
		//VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id;`,
		//	anime.Name, anime.Russian, anime.English, anime.Japanese, anime.Synonyms, anime.Status,
		//	anime.Episodes, anime.EpisodesAired, anime.AiredOn, anime.ReleasedOn, anime.Duration, nil, nil,
		//	anime.UpdatedAt, anime.NextEpisodeAt, nil, nil)

		//var id string
		//err = row.Scan(&id)
		//if err != nil {
		//	log.Printf("Error when saving anime %d: %q", aid, err.Error())
		//	continue
		//}
		//log.Printf("Anime %q saved!", anime.Name)

	}

	for i := 1; i < maxNumberPage; i++ {
		person, err := shikimori_api.GetPersonInfo(i)
		if err != nil {
			log.Printf("Error in person %d: %q", i, err.Error())
		} else {
			log.Printf("Person %d, name %q exists", i, person.Name)
		}
	}
}
