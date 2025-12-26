package main

import (
	"context"
	"log"

	"github.com/qqMelon/restorely-api/internal/restore"
)

func main () {
	// Test on backup.dump at root project
	err := restore.TestPostgresRestoreDocker(
		context.Background(),
		"./backup.dump",
	)

	if err != nil {
		log.Fatalf("RESTORE FAILED: %v", err)
	}

	log.Println("RESTORE OK")
}
