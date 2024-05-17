package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type ImageData struct {
	Width        int    `json:"width,omitempty"` //`json:"-"`
	Height       int    `json:"height,omitempty"`
	Format       string `json:"format,omitempty"`
	FormatSource string `json:"formatSource,omitempty"`
}

type Image struct {
	Path      string
	SourseImg string
	Meta      ImageData
}

func GetImage(sourceURL string) *Image {
	dbImg := &Image{}

	row := Conn.QueryRow(context.Background(), `
		SELECT path, source_img, meta FROM public.images_table WHERE source_img=$1`, sourceURL)

	err := row.Scan(&dbImg.Path, &dbImg.SourseImg, &dbImg.Meta)
	if errors.Is(pgx.ErrNoRows, err) {
		return nil
	}
	if err != nil {
		log.Printf("Error in get image: %q", err.Error())
	}

	return dbImg
}

func CreateOrUpdateImage(path string, sourceUrl string, meta ImageData) error {
	_, err := Conn.Exec(context.Background(), `
		INSERT INTO public.images_table (path, source_img, meta) VALUES ($1, $2, $3)
		ON CONFLICT (source_img) DO UPDATE
		SET meta = $3`, path, sourceUrl, meta,
	)
	if err != nil {
		log.Printf("Error in update image: %q", err.Error())
		return err
	}
	fmt.Printf("Add '%s' in Image table\n", path)
	return nil
}
