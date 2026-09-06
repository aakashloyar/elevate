package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in"
)

type Handler struct{ service in.GetAttemptProblemsService }

func NewHandler(service in.GetAttemptProblemsService) *Handler { return &Handler{service: service} }

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) GetAttemptProblems(w http.ResponseWriter, r *http.Request, attemptID string) {
	offset, limit, err := pagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := h.service.Execute(r.Context(), in.GetAttemptProblemsInput{AttemptID: attemptID, Offset: offset, Limit: limit})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func pagination(r *http.Request) (int, int, error) {
	offset, err := integerQuery(r, "offset", 0)
	if err != nil || offset < 0 {
		return 0, 0, &queryError{"offset must be a non-negative integer"}
	}
	limit, err := integerQuery(r, "limit", 10)
	if err != nil || limit <= 0 {
		return 0, 0, &queryError{"limit must be a positive integer"}
	}
	return offset, limit, nil
}
func integerQuery(r *http.Request, name string, fallback int) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}

type queryError struct{ message string }

func (e *queryError) Error() string { return e.message }
