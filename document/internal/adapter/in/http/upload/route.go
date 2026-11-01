package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/documents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ListDocuments(w, r)
	})
	mux.HandleFunc("/documents/batch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.GetDocumentsBatch(w, r)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.InitUpload(w, r)
			}
		default:
			{
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		}
	})
	mux.HandleFunc("/complete-upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				h.CompleteUpload(w, r)
			}
		default:
			{
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		}
	})

}
