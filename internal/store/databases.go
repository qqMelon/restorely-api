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

func GetDatabaseByID(
	ctx context.Context,
	db *sql.DB,
	id string,
)(*models.Database, error) {
	row := db.QueryRowContext(ctx, `
		SELECT
			d.id,
			d.name,
			d.type,
			d.host,
			d.port,
			d.db_name,
			b.status,
			b.created_at,
			b.duration_ms,
			r.status,
			r.created_at,
			r.duration_ms
		FROM databases d
		LEFT JOIN LATERAL (
			SELECT status, created_at, duration_ms
			FROM backups
			WHERE database_id = d.id
			ORDER BY created_at DESC
			LIMIT 1
		) b ON true
		LEFT JOIN LATERAL (
			SELECT rt.status, rt.created_at, rt.duration_ms
			FROM restore_tests rt
			JOIN backups b2 ON b2.id = rt.backup_id
			WHERE b2.database_id = d.id
			ORDER BY rt.created_at DESC
			LIMIT 1
		) r ON true
		WHERE d.id = $1
	`, id)

	var (
		d models.Database

		bStatus		sql.NullString
		bCreated 	sql.NullTime
		bDuration sql.NullInt64

		rStatus		sql.NullString
		rCreated	sql.NullTime
		rDuration sql.NullInt64
	)

	err := row.Scan(
		&d.ID,
		&d.Name,
		&d.Type,
		&d.Host,
		&d.Port,
		&d.DBName,
		&bStatus,
		&bCreated,
		&bDuration,
		&rStatus,
		&rCreated,
		&rDuration,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if bStatus.Valid {
		d.LastBackup = &models.Backup{
			Status: bStatus.String,
			CreatedAt: bCreated.Time,
			DurationMS: bDuration.Int64,
		}
	}

	if rStatus.Valid {
		d.LastRestoreTest = &models.RestoreTest{
			Status: rStatus.String,
			CreatedAt: rCreated.Time,
			DurationMS: rDuration.Int64,
		}
	}

	return &d, nil
}
