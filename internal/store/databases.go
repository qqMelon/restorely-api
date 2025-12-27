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
		  SELECT *
		  FROM backups
		  WHERE database_id = d.id
		  ORDER BY created_at DESC
		  LIMIT 1
		) b ON true
		LEFT JOIN LATERAL (
		  SELECT *
		  FROM restore_tests
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

		var (
			bStatus   sql.NullString
			bCreated  sql.NullTime
			bDuration sql.NullInt64

			rStatus   sql.NullString
			rCreated  sql.NullTime
			rDuration sql.NullInt64
		)

		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Type,
			&bStatus,
			&bCreated,
			&bDuration,
			&rStatus,
			&rCreated,
			&rDuration,
		)
		if err != nil {
			return nil, err
		}

		if bStatus.Valid {
			d.LastBackup = &models.Backup{
				Status:     bStatus.String,
				CreatedAt:  bCreated.Time,
				DurationMS: bDuration.Int64,
			}
		} else {
			d.LastBackup = nil
		}

		if rStatus.Valid {
			d.LastRestoreTest = &models.RestoreTest{
				Status:     rStatus.String,
				CreatedAt:  rCreated.Time,
				DurationMS: rDuration.Int64,
			}
		} else {
			d.LastRestoreTest = nil
		}

		result = append(result, d)
	}

	return result, nil
}

