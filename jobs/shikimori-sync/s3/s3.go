package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/google/uuid"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3api "github.com/aws/aws-sdk-go/service/s3"
)

var s3Conn *s3api.S3
var s3Bucket string
var s3Session *session.Session

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
		WithLogLevel(aws.LogOff).
		WithRegion(region).
		WithEndpoint("s3." + region + ".scw.cloud")
	s3Conn = s3api.New(s3Session, cfg)
}

func getFileHashFromS3(s3Url string) string {

	headObj, err := s3Conn.HeadObject(&s3api.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Url),
	})

	if err != nil {
		log.Fatalf("Error in get hash object: %q", err.Error())
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
