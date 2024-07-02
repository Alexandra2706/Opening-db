package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

var existGenres map[int]string = make(map[int]string)
var existGenresLock sync.Mutex

func CreateOrUpdateGenre(shikimoriId int, genreName string, russian string) string {
	existGenresLock.Lock()
	existId, ok := existGenres[shikimoriId]
	existGenresLock.Unlock()
	if ok {
		return existId
	}

	tempStr := strings.ToLower(genreName) // заменить пробелы на нижн подчеркивание
	stringId := strings.Replace(tempStr, " ", "_", -1)

	err := Conn.QueryRow(context.Background(), `
		INSERT INTO public.genres_table (id, shikimori_id, genre_name, russian) VALUES ($1, $2, $3, $4)
		ON CONFLICT (shikimori_id) DO UPDATE
		SET genre_name = $3, russian = $4 RETURNING id`, stringId, shikimoriId, genreName, russian).Scan(&stringId)
	if err != nil {
		log.Fatalf("Error in update genre: %q", err.Error())
	}
	fmt.Printf("Add '%s' in Genre table\n", stringId)

	existGenresLock.Lock()
	existGenres[shikimoriId] = stringId
	existGenresLock.Unlock()

	return stringId
}
