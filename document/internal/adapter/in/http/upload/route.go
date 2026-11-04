package http

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/documents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ListDocuments(w, r)
	})
	mux.HandleFunc("/documents/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/documents/")
		if r.Method != http.MethodGet || path == "" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if strings.HasSuffix(path, "/content") {
			documentID := strings.TrimSuffix(path, "/content")
			if documentID == "" || strings.Contains(documentID, "/") {
				http.NotFound(w, r)
				return
			}
			h.GetDocumentContent(w, r, documentID)
			return
		}
		if strings.Contains(path, "/") {
			http.NotFound(w, r)
			return
		}
		h.GetDocument(w, r, path)
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
