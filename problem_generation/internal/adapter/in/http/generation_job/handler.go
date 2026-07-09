package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in/generation_job"
)

type CreateGenerationJobRequest struct {
	UserID             string   `json:"user_id"`
	SingleCorrectCount int      `json:"single_correct_count"`
	MultiCorrectCount  int      `json:"multi_correct_count"`
	NumericalCount     int      `json:"numerical_count"`
	DocumentID         *string  `json:"document_id"`
	AssessmentID       *string  `json:"assessment_id"`
	Level              string   `json:"level"`
	Description        string   `json:"description"`
	TopicIDs           []string `json:"topic_ids"`
}

type CreateGenerationJobResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

type GetGenerationJobResponse struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	SingleCorrectCount int      `json:"single_correct_count"`
	MultiCorrectCount  int      `json:"multi_correct_count"`
	NumericalCount     int      `json:"numerical_count"`
	DocumentID         *string  `json:"document_id"`
	AssessmentID       *string  `json:"assessment_id"`
	Level              string   `json:"level"`
	Description        string   `json:"description"`
	Status             string   `json:"status"`
	TopicIDs           []string `json:"topic_ids"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

type Handler struct {
	createGenerationJobService in.CreateGenerationJobService
	getGenerationJobService    in.GetGenerationJobService
}

func NewHandler(createGenerationJobService in.CreateGenerationJobService, getGenerationJobService in.GetGenerationJobService) *Handler {
	return &Handler{
		createGenerationJobService: createGenerationJobService,
		getGenerationJobService:    getGenerationJobService,
	}
}

func (h *Handler) CreateGenerationJob(w http.ResponseWriter, r *http.Request) {
	var req CreateGenerationJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.createGenerationJobService.Execute(r.Context(), in.CreateGenerationJobInput{
		UserID:             req.UserID,
		SingleCorrectCount: req.SingleCorrectCount,
		MultiCorrectCount:  req.MultiCorrectCount,
		NumericalCount:     req.NumericalCount,
		DocumentID:         req.DocumentID,
		AssessmentID:       req.AssessmentID,
		Level:              req.Level,
		Description:        req.Description,
		TopicIDs:           req.TopicIDs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateGenerationJobResponse{JobID: out.JobID, Status: out.Status})
}

func (h *Handler) GetGenerationJobByID(w http.ResponseWriter, r *http.Request, jobID string) {
	out, err := h.getGenerationJobService.Execute(r.Context(), in.GetGenerationJobInput{JobID: jobID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetGenerationJobResponse{
		ID:                 out.ID,
		UserID:             out.UserID,
		SingleCorrectCount: out.SingleCorrectCount,
		MultiCorrectCount:  out.MultiCorrectCount,
		NumericalCount:     out.NumericalCount,
		DocumentID:         out.DocumentID,
		AssessmentID:       out.AssessmentID,
		Level:              out.Level,
		Description:        out.Description,
		Status:             out.Status,
		TopicIDs:           out.TopicIDs,
		CreatedAt:          out.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:          out.UpdatedAt.Format(http.TimeFormat),
	})
}

func (h *Handler) IsGenerationJobRoute(path string) bool {
	return strings.HasPrefix(path, "/generation-jobs")
}
