package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/notifications/email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.SendEmail(w, r)
	})
}
