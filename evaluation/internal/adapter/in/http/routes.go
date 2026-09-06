package httpapi

import (
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/evaluation/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Health(w, r)
	})

	mux.HandleFunc("/evaluations/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		submissionID := strings.TrimPrefix(r.URL.Path, "/evaluations/")
		if submissionID == "" || strings.Contains(submissionID, "/") {
			http.NotFound(w, r)
			return
		}

		h.GetEvaluationBySubmissionID(w, r, submissionID)
	})
}
