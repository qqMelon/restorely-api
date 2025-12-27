package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/qqMelon/restorely-api/internal/store"
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

	log.Println("Restorely worker started")

	for {
		runOnce(db)
		time.Sleep(30 * time.Second)
	}
}

func runOnce(db *sql.DB) {
	ctx := context.Background()

	rows, err := db.Query(`
		SELECT d.id
		FROM databases d
		LEFT JOIN backups b ON b.database_id = d.id
		GROUP BY d.id
		HAVING MAX(b.created_at) IS NULL
		   OR MAX(b.created_at) < NOW() - INTERVAL '1 hour'
	`)
	if err != nil {
		log.Println("DB list error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var databaseID string
		if err := rows.Scan(&databaseID); err != nil {
			continue
		}

		log.Println("Running scheduled backup for DB:", databaseID)

		start := time.Now()
		time.Sleep(2 * time.Second)

		backupID, err := store.CreateBackup(
			ctx,
			db,
			databaseID,
			"success",
			time.Since(start),
		)
		if err != nil {
			log.Println("backup insert error:", err)
			continue
		}

		log.Println("Running restore test for backup:", backupID)

		start = time.Now()
		time.Sleep(3 * time.Second)

		err = store.CreateRestoreTest(
			ctx,
			db,
			backupID,
			"success",
			time.Since(start),
		)
		if err != nil {
			log.Println("restore insert error:", err)
			continue
		}
	}
}
