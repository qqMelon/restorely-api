package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateBackup (
	ctx context.Context,
	db *sql.DB,
	databaseID string,
	status string,
	duration time.Duration,
)(string, error) {
	id := uuid.New().String()

	_, err := db.ExecContext(ctx,`
		INSERT INTO backups (id, database_id, status, duration_ms)
		VALUES ($1, $2, $3, $4)
	`,
		id,
		databaseID,
		status,
		duration.Milliseconds(),
	)

	return id, err
}

func CreateRestoreTest (
	ctx context.Context,
	db *sql.DB,
	backupID string,
	status string,
	duration time.Duration,
) error {
	id := uuid.New().String()

	_, err := db.ExecContext(ctx, `
		INSERT INTO restore_tests (id, backup_id, status, duration_ms)
		VALUES ($1, $2, $3, $4)
	`,
		id, 
		backupID,
		status,
		duration.Milliseconds(),
	)

	return err
}
