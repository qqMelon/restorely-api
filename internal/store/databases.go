package store

import (
	"context"
	"database/sql"

	"github.com/qqMelon/restorely-api/internal/models"
)

func ListDatabases(ctx context.Context, db *sql.DB) ([]models.Database, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT
		  d.id, d.name, d.type,
		  b.status, b.created_at, b.duration_ms,
		  r.status, r.created_at, r.duration_ms
		FROM databases d
		LEFT JOIN LATERAL (
		  SELECT * FROM backups
		  WHERE database_id = d.id
		  ORDER BY created_at DESC
		  LIMIT 1
		) b ON true
		LEFT JOIN LATERAL (
		  SELECT * FROM restore_tests
		  WHERE backup_id = b.id
		  ORDER BY created_at DESC
		  LIMIT 1
		) r ON true
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Database

	for rows.Next() {
		var d models.Database
		var b models.Backup
		var r models.RestoreTest

		err := rows.Scan(
			&d.ID, &d.Name, &d.Type,
			&b.Status, &b.CreatedAt, &b.DurationMS,
			&r.Status, &r.CreatedAt, &r.DurationMS,
		)
		if err != nil {
			return nil, err
		}

		d.LastBackup = &b
		d.LastRestoreTest = &r

		result = append(result, d)
	}

	return result, nil
}
