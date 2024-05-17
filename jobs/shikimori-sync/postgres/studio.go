package postgres

import (
	"context"
	"fmt"
	"log"
	//"strings"
	"sync"

	"github.com/google/uuid"
)

const (
	imageEndpoint = "https://shikimori.one"
)

var existStudios map[int]uuid.UUID = make(map[int]uuid.UUID)
var existStudiosLock sync.Mutex

func CreateOrUpdateStudio(shikimoriId int, studioName string, imageUrl string) uuid.UUID {
	existStudiosLock.Lock()
	existId, ok := existStudios[shikimoriId]
	existStudiosLock.Unlock()
	if ok {
		return existId
	}

	uuidId := uuid.New()

	fmt.Println(uuidId, shikimoriId, studioName, imageUrl)
	//
	//path, err := s3.ImportImage(imageEndpoint + imageUrl)
	//if err != nil {
	//	log.Printf("Error in get image: %q", err.Error())
	//}

	err := Conn.QueryRow(context.Background(), `
		INSERT INTO public.studio_table (id, shikimori_id, studio_name, image) VALUES ($1, $2, $3, $4)
		ON CONFLICT (shikimori_id) DO UPDATE
		SET studio_name = $3, image = $4 RETURNING id`, uuidId, shikimoriId, studioName, imageUrl,
	).Scan(&uuidId)
	if err != nil {
		log.Printf("Error in update studio: %q", err.Error())
	}
	fmt.Printf("Add '%s' in Studio table\n", studioName)

	existStudiosLock.Lock()
	existStudios[shikimoriId] = uuidId
	existStudiosLock.Unlock()

	return uuidId
}
