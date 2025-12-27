package httpapi

import "net/http"

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/databases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.CreateDatabase(w, r)
			return 
		}

		if r.Method == http.MethodGet {
			s.ListDatabases(w, r)
			return 
		}
	})

	return withCORS(mux)
}
