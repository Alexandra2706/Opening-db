package main

//go:generate swag init

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/swaggo/http-swagger/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	_ "api/docs"
	"api/postgres"
	"api/routes"
)

// @title OPDB Test API
// @version 1.0
// @description This is a OPDB Test server.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1

func main() {

	defer postgres.CloseConnection()

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	routes.Init(mux)

	//h2cWrapper := &h2c.HandlerH2C{
	//	Handler:  cors.Default().Handler(mux),
	//	H2Server: &http2.Server{},
	//}

	http2Server := &http2.Server{}

	corsPolicy := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:      []string{"*"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(corsPolicy.Handler(mux), http2Server),
	}

	err := http2.ConfigureServer(server, &http2.Server{})
	if err != nil {
		log.Fatal("http2.ConfigureServer: ", err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
