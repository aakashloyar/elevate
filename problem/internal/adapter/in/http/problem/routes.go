package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/problems", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListProblems(w, r)
		case http.MethodPost:
			h.CreateProblem(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/problems/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/problems/")
		parts := strings.Split(path, "/")
		problemID := parts[0]
		if problemID == "" {
			http.Error(w, "missing problem id", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				h.GetProblemByID(w, r, problemID)
			case http.MethodPut:
				h.UpdateProblem(w, r, problemID)
			case http.MethodDelete:
				h.DeleteProblem(w, r, problemID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}
