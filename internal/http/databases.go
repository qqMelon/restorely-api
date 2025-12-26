package httpapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qqMelon/restorely-api/internal/store"
)

type Server struct {
	DB *sql.DB
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
