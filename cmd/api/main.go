package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpapi "github.com/qqMelon/restorely-api/internal/http"
)

func main () {
	db, err := sql.Open("postgres", "postgres://restorely:restorely@localhost:5433/restorely?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := &httpapi.Server{DB: db}

	log.Println("API Running on :8080")
	http.ListenAndServe(":8080", server.Router())
}
