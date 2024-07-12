package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func init() { //выполняется всегда до main.go
	var err error
	Conn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %q\n", err)
	}
}

func CloseConnection() {
	err := Conn.Close(context.Background())
	if err != nil {
		log.Fatalf("Unable to close connection to database: %q\n", err)
	}
}
