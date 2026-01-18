package httpapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/qqMelon/restorely-api/internal/store"
)

type HistoryItems struct {
	Type string `json:"type"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
	DurationMS int64 `json:"duration_ms"`
}

type CreateDatabaseRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName string `json:"db_name"`
}

type Server struct {
	DB *sql.DB
}

func (s *Server) DatabaseHistory(w http.ResponseWriter, r *http.Request) {
	dbID := chi.URLParam(r, "id")

	rows, err := s.DB.Query(`
		SELECT
		  'backup' AS type,
		  b.status,
		  b.created_at,
		  b.duration_ms
		FROM backups b
		WHERE b.database_id = $1

		UNION ALL

		SELECT
		  'restore' AS type,
		  rt.status,
		  rt.created_at,
		  rt.duration_ms
		FROM restore_tests rt
		JOIN backups b ON b.id = rt.backup_id
		WHERE b.database_id = $1

		ORDER BY created_at DESC
		LIMIT 50
	`, dbID)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var history []HistoryItems

	for rows.Next() {
		var h HistoryItems
		if err := rows.Scan(
			&h.Type,
			&h.Status,
			&h.CreatedAt,
			&h.DurationMS,
		); err != nil {
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		history = append(history, h)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (s *Server) ListDatabases(w http.ResponseWriter, r *http.Request) {
	databases, err := store.ListDatabases(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "internal error", 500)
		fmt.Println("Error on databases getter:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}

func (s *Server) GetDatabase(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid database id", http.StatusBadRequest)
		return
	}

	db, err := store.GetDatabaseByID(r.Context(), s.DB, id)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}

	if db == nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(db)
}

func (s *Server) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDatabase called")
	var req CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", 400)
		return
	}

	dsn := "postgres://" + req.Username + ":" + req.Password + "@" + req.Host + ":" + strconv.Itoa(req.Port) + "/" + req.DBName + "?sslmode=disable"

	testDB, err := sql.Open("postgres", dsn)
	if err != nil {
		http.Error(w, "connection failed", 400)
		return
	}
	defer testDB.Close()

	if err := testDB.Ping(); err != nil {
		http.Error(w, "cannot connect to database", 400)
		return
	}

	id := uuid.New().String()

	_, err = s.DB.Exec(`
		INSERT INTO databases
		(id, name, type, host, port, username, password, db_name)
		VALUES ($1, $2, 'postgres', $3, $4, $5, $6, $7)
	`,
		id,
		req.Name,
		req.Host,
		req.Port,
		req.Username,
		req.Password,
		req.DBName,
	)
	if err != nil {
		http.Error(w, "insert failed", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
