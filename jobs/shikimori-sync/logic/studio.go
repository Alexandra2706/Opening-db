package logic

import (
	"github.com/google/uuid"
	"log"
	"shikimori-sync/postgres"
	"shikimori-sync/s3"
)

func CreateOrUpdateStudio(studioID int, studioName string, imageUrl string) uuid.UUID {
	err := s3.CreateOrUpdateImage(imageUrl)
	if err != nil {
		log.Fatalf("Error in get image: %q", err.Error())
	}
	log.Println("Studio data: ", studioID, studioName, imageUrl)
	return postgres.CreateOrUpdateStudio(studioID, studioName, imageUrl)
}
