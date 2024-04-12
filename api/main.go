package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	fmt.Println("Hello world")
	http.HandleFunc("/", sayhello)

	server := &http.Server{
		Addr: ":8080",
	}
	http2.ConfigureServer(server, &http2.Server{})

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет!")
}
