package httpapi

import "net/http"

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/databases", s.ListDatabases)

	return withCORS(mux)
}
