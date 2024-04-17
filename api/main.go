package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", sayhello)

	//h2cWrapper := &h2c.HandlerH2C{
	//	Handler:  cors.Default().Handler(mux),
	//	H2Server: &http2.Server{},
	//}

	http2Server := &http2.Server{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(cors.Default().Handler(mux), http2Server),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World http2"))
}
