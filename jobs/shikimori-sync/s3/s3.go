package s3

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/google/uuid"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"shikimori-sync/postgres"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3api "github.com/aws/aws-sdk-go/service/s3"
)

var s3Conn *s3api.S3
var s3Bucket string
var s3Session *session.Session

func CreateOrUpdateImage(sourceUrl string) error {
	// Скачать картинку с помощью get запроса
	fmt.Println(sourceUrl)
	resp, err := http.Get(sourceUrl)
	if err != nil {
		log.Fatalln(err)
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

	fmt.Println(mimeType)

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

	fmt.Println("type: ", reflect.TypeOf(originalImage))
	imgWidth := originalImage.Bounds().Dx()
	imgHeight := originalImage.Bounds().Dy()
	fmt.Println("Image width: ", imgWidth)
	fmt.Println("Image height: ", imgHeight)

	//Посчитать md5 - хэш сумму конвертированной картинки

	hash := md5.Sum(imgBytes)
	fmt.Println("hash: ", hash)
	checkSum := hex.EncodeToString(hash[:])
	fmt.Println("checkSum: ", checkSum)

	//Проверить в таблице image_table по source_url есть ли запись
	dbImg := &postgres.Image{}

	dbImg = postgres.GetImage(sourceUrl)

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
				return err
			}
			DeleteFileFromS3(dbImg.Path)
			return nil
		} else if getHashImgFromS3(dbImg.Path) != checkSum {
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
		return err
	}

	//Да - сверяем md5 и mime type, если сходится, то возвращаем path, нет - перезаписываем файл в s3 и возвращаем path,
	//	 - нет и поменялся mime type: загружаем новую картинку с новым именем, обновляем запись в БД, удаляем старую записьв БД
	//Нет
	//    - генерируем имя для s3 images/*[uuid].png или jpeg
	//    - сохранием файл на s3 (указать правильный mime type)
	// 	  - добавляем запись в image_table
	//    - при ошибке записи в таблицу, удалить файл из s3

	return nil
}

type simpleLogger struct{}

func (*simpleLogger) Log(a ...interface{}) { log.Println(a...) }

func init() {
	s3Bucket = os.Getenv("S3_BUCKET")
	s3Session = session.Must(session.NewSession())
	region := os.Getenv("S3_REGION")
	if region == "" {
		region = "fr-par"
	}
	cfg := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(
		os.Getenv("S3_PUBLIC_KEY"), os.Getenv("S3_PRIVATE_KEY"), "")).
		WithMaxRetries(1).
		WithLogger(&simpleLogger{}).
		WithLogLevel(aws.LogDebug).
		WithRegion(region).
		WithEndpoint("s3." + region + ".scw.cloud")
	s3Conn = s3api.New(s3Session, cfg)
}

func getHashImgFromS3(s3Url string) string {

	headObj, err := s3Conn.HeadObject(&s3api.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Url),
	})

	if err != nil {
		log.Printf("Error in get hash object: %q", err.Error())
	}
	return *headObj.ETag
}

// Создать новую запись в s3
func UploadFileToS3(imgBytes []byte, mimeType string) string {

	//создать путь
	ext := strings.Split(mimeType, "/")
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	path := "images/" + newUUID.String() + "." + ext[1]

	reader := bytes.NewReader(imgBytes)
	_, err = s3Conn.PutObject(&s3api.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
		Body:   reader,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully uploaded to ", path)
	return path
}

// обновить файл на s3
func UpdateFileOnS3(imgBytes []byte, path string) string {

	reader := bytes.NewReader(imgBytes)
	_, err := s3Conn.PutObject(&s3api.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
		Body:   reader,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully updated on ", path)
	return path
}

// Удалить запись из s3
func DeleteFileFromS3(path string) {
	object, err := s3Conn.DeleteObject(&s3api.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully deleted object ", object)
}
