package s3

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"reflect"
	"shikimori-sync/postgres"
)

func CreateOrUpdateImage(sourceUrl string) error {
	// Скачать картинку с помощью get запроса
	log.Println(sourceUrl)
	resp, err := http.Get(sourceUrl)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code: %d", resp.StatusCode))
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot read body: %q", err.Error()))
	}

	// Посмотреть mime type изображения
	mimeType := http.DetectContentType(imgBytes)

	log.Println(mimeType)

	//По mime type определить надо ли конвертировать изображение и как
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return errors.New(fmt.Sprintf("Unsuppoted image mime type: %q", mimeType))
		// тут надо сделать конвертицию изображения в jpg
	}

	//Конвертация
	//открыть с помощью пакета image
	//вытащить из нее размер и проверить, что она открылась без ошибки

	originalImage, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {

		return errors.New(fmt.Sprintf("failed to decode image: %q", err.Error()))
	}

	log.Println("type: ", reflect.TypeOf(originalImage))
	imgWidth := originalImage.Bounds().Dx()
	imgHeight := originalImage.Bounds().Dy()
	log.Println("Image width: ", imgWidth)
	log.Println("Image height: ", imgHeight)

	//Посчитать md5 - хэш сумму конвертированной картинки

	hash := md5.Sum(imgBytes)
	log.Println("hash: ", hash)
	checkSum := hex.EncodeToString(hash[:])
	log.Println("checkSum: ", checkSum)

	dbImg := postgres.GetImage(sourceUrl)

	// если запись в бд существует
	if dbImg != nil {
		//проверяем mime-type и md5
		if dbImg.Meta.Format != mimeType {
			newPath := UploadFileToS3(imgBytes, mimeType)
			err := postgres.CreateOrUpdateImage(newPath, sourceUrl, postgres.ImageData{
				Width:        imgWidth,
				Height:       imgHeight,
				Format:       mimeType,
				FormatSource: mimeType,
			})
			if err != nil {
				DeleteFileFromS3(newPath)
				return err
			}
			DeleteFileFromS3(dbImg.Path)
		} else if getFileHashFromS3(dbImg.Path) != checkSum {
			UpdateFileOnS3(imgBytes, dbImg.Path)
		}
		return nil
	}

	newPath := UploadFileToS3(imgBytes, mimeType)
	err = postgres.CreateOrUpdateImage(newPath, sourceUrl, postgres.ImageData{
		Width:        imgWidth,
		Height:       imgHeight,
		Format:       mimeType,
		FormatSource: mimeType,
	})
	if err != nil {
		DeleteFileFromS3(newPath)
		return err
	}

	return nil
}
