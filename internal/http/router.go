package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/databases", func(r chi.Router) {
		r.Get("/", s.ListDatabases)
		r.Post("/", s.CreateDatabase)
		r.Get("/{id}", s.GetDatabase)
		r.Get("/{id}/history", s.DatabaseHistory)
	})

	return withCORS(r)
}
