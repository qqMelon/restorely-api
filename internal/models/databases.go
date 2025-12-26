package models

import "time"

type Database struct {
	ID	string 	`json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	LastBackup *Backup `json:"last_backup"`
	LastRestoreTest *RestoreTest `json:"last_restore_test"`
}

type Backup struct {
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	DurationMS int64 `json:"duration_ms"`
}

type RestoreTest struct {
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	DurationMS int64 `json:"duration_ms"`
}
