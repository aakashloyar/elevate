package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("welcome to user service of elevate"))
		case http.MethodPost:
			h.CreateUser(w, r)	
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		parts := strings.Split(path, "/")
		userID := parts[0]
		if userID == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				h.GetUserByID(w, r, userID)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})
}
