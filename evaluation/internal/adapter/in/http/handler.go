package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	in "github.com/aakashloyar/elevate/evaluation/internal/application/ports/in"
	evaluationservice "github.com/aakashloyar/elevate/evaluation/internal/application/service"
)

type Handler struct {
	getEvaluationService in.GetEvaluationService
}

func NewHandler(getEvaluationService in.GetEvaluationService) *Handler {
	return &Handler{getEvaluationService: getEvaluationService}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) GetEvaluationBySubmissionID(w http.ResponseWriter, r *http.Request, submissionID string) {
	evaluation, err := h.getEvaluationService.Execute(r.Context(), submissionID)
	if errors.Is(err, evaluationservice.ErrEvaluationNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to load evaluation", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(evaluation)
}
