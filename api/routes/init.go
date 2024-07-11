package routes

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/anime", listAnime)
	mux.HandleFunc("GET /v1/person", listPerson)
}
