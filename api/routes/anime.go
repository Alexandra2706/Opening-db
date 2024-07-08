package routes

import (
	"api/postgres"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type errorResponse struct {
	Message string `json:"message"`
}

// @Summary Get Anime List
// @Tags Lists
// @Description Get anime list
// @ID get-anime-list
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} []structures.AnimeShortInfo
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /anime [get]
func listAnime(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 1000
	var err error

	//json.NewDecoder(r.Body).Decode(&offset) Оставить пока
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage, err := json.Marshal(errorResponse{Message: "Offset must be int"})
			if err == nil {
				w.Write(errMessage)
			}
			return
		}
	}
	if offset < 0 {
		offset = 0
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMessage, err := json.Marshal(errorResponse{Message: "Limit must be int"})
			if err == nil {
				w.Write(errMessage)
			}
			return
		}
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	//http.StatusBadRequest, err.Error())

	listAnimeInfo, err := postgres.ListAnime(limit, offset)
	if err != nil {
		log.Println("get list anime ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errMessage, err := json.Marshal(errorResponse{Message: "Internal Server Error"})
		if err == nil {
			w.Write(errMessage)
		}
		return
	}
	jsonListAnime, err := json.MarshalIndent(&listAnimeInfo, "", " ")
	if err != nil {
		log.Println("marshal ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("JSON string: %s\n", string(jsonListAnime))

	fmt.Println("Ofsset=", offsetStr)
	fmt.Println("Limit=", limitStr)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonListAnime)
	date := time.Date(2020, 9, 9, 13, 34, 17, 0, time.UTC)
	fmt.Println("date=", date)
}
