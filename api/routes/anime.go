package routes

import (
	"api/postgres"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type errorResponse struct {
	Message string `json:"message"`
}

func PrepareLimitOffset(w http.ResponseWriter, r *http.Request) (*int, *int, error) {
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
			return nil, nil, err
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
			return nil, nil, err
		}
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	return &limit, &offset, nil
}

// @Summary Get Anime List
// @Tags Lists
// @Description Get anime list
// @ID get-anime-list
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} []structures.AnimeInfo
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /anime [get]
func listAnime(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := PrepareLimitOffset(w, r)
	if err != nil {
		return
	}
	listAnimeInfo, err := postgres.ListAnime(*limit, *offset)
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
	log.Printf("JSON string: %s\n", string(jsonListAnime))

	w.WriteHeader(http.StatusOK)
	w.Write(jsonListAnime)
}

// @Summary Get person List
// @Tags Lists
// @Description Get person list
// @ID get-person-list
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {object} []structures.PersonInfo
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /person [get]
func listPerson(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := PrepareLimitOffset(w, r)
	if err != nil {
		return
	}
	listPersonInfo, err := postgres.ListPerson(*limit, *offset)
	if err != nil {
		log.Println("get list anime ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errMessage, err := json.Marshal(errorResponse{Message: "Internal Server Error"})
		if err == nil {
			w.Write(errMessage)
		}
		return
	}
	jsonListPerson, err := json.MarshalIndent(&listPersonInfo, "", " ")
	if err != nil {
		log.Println("marshal ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("JSON string: %s\n", string(jsonListPerson))

	w.WriteHeader(http.StatusOK)
	w.Write(jsonListPerson)

}
