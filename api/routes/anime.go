package routes

import "net/http"

func listAnime(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("listAnime"))
}
