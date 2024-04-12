package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", sayhello)

	server := &http.Server{
		Addr:    ":8080",
		Handler: cors.Default().Handler(mux),
	}
	http2.ConfigureServer(server, &http2.Server{})

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World http2"))
}
