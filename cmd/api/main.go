package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	httpapi "github.com/qqMelon/restorely-api/internal/http"
)

func main () {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL missing")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	server := &httpapi.Server{DB: db}

	log.Println("API Running on :8080")
	http.ListenAndServe(":8080", server.Router())
}
