package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/evaluation/internal/application"
)

type Handler struct{ service *application.Service }

func New(service *application.Service) *Handler { return &Handler{service: service} }
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/evaluation/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_, _ = w.Write([]byte("ok"))
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
		evaluation, err := h.service.Get(r.Context(), submissionID)
		if errors.Is(err, application.ErrEvaluationNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "failed to load evaluation", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(evaluation)
	})
}
