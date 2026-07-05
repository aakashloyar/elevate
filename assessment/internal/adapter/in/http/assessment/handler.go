package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in/assessment"
)

type CreateAssessmentRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	DurationSeconds int    `json:"duration_seconds"`
	CreatedBy       string `json:"created_by"`
}

type CreateAssessmentResponse struct {
	AssessmentID string `json:"assessment_id"`
}

type GetAssessmentResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	DurationSeconds int    `json:"duration_seconds"`
	CreatedBy       string `json:"created_by"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type Handler struct {
	createAssessmentService in.CreateAssessmentService
	getAssessmentService    in.GetAssessmentService
	deleteAssessmentService in.DeleteAssessmentService
}

func NewHandler(createAssessmentService in.CreateAssessmentService, getAssessmentService in.GetAssessmentService, deleteAssessmentService in.DeleteAssessmentService) *Handler {
	return &Handler{
		createAssessmentService: createAssessmentService,
		getAssessmentService:    getAssessmentService,
		deleteAssessmentService: deleteAssessmentService,
	}
}

func (h *Handler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
	var req CreateAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.createAssessmentService.Execute(r.Context(), in.CreateAssessmentInput{
		Title:           req.Title,
		Description:     req.Description,
		Status:          req.Status,
		DurationSeconds: req.DurationSeconds,
		CreatedBy:       req.CreatedBy,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateAssessmentResponse{AssessmentID: out.AssessmentID})
}

func (h *Handler) GetAssessmentByID(w http.ResponseWriter, r *http.Request, assessmentID string) {
	out, err := h.getAssessmentService.Execute(r.Context(), in.GetAssessmentInput{AssessmentID: assessmentID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetAssessmentResponse{
		ID:              out.ID,
		Title:           out.Title,
		Description:     out.Description,
		Status:          out.Status,
		DurationSeconds: out.DurationSeconds,
		CreatedBy:       out.CreatedBy,
		CreatedAt:       out.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:       out.UpdatedAt.Format(http.TimeFormat),
	})
}

func (h *Handler) DeleteAssessment(w http.ResponseWriter, r *http.Request, assessmentID string) {
	if err := h.deleteAssessmentService.Execute(r.Context(), in.DeleteAssessmentInput{AssessmentID: assessmentID}); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) IsAssessmentRoute(path string) bool {
	return strings.HasPrefix(path, "/assessments")
}
