package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
)

type CreateAssessmentRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
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
	DurationSeconds int    `json:"duration_seconds"`
	CreatedBy       string `json:"created_by"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type Handler struct {
	createAssessmentService in.CreateAssessmentService
	listAssessmentsService  in.ListAssessmentsService
	getAssessmentService    in.GetAssessmentService
	deleteAssessmentService in.DeleteAssessmentService
}

func NewHandler(createAssessmentService in.CreateAssessmentService, listAssessmentsService in.ListAssessmentsService, getAssessmentService in.GetAssessmentService, deleteAssessmentService in.DeleteAssessmentService) *Handler {
	return &Handler{
		createAssessmentService: createAssessmentService,
		listAssessmentsService:  listAssessmentsService,
		getAssessmentService:    getAssessmentService,
		deleteAssessmentService: deleteAssessmentService,
	}
}

type ListAssessmentsResponse struct {
	Assessments []GetAssessmentResponse `json:"assessments"`
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

func (h *Handler) ListAssessments(w http.ResponseWriter, r *http.Request) {
	input := in.ListAssessmentsInput{}

	if userID := strings.TrimSpace(r.URL.Query().Get("userID")); userID != "" {
		input.UserID = &userID
	}
	if title := strings.TrimSpace(r.URL.Query().Get("title")); title != "" {
		input.Title = &title
	}
	if description := strings.TrimSpace(r.URL.Query().Get("description")); description != "" {
		input.Description = &description
	}
	if limit := strings.TrimSpace(r.URL.Query().Get("limit")); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil || parsedLimit < 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		input.Limit = &parsedLimit
	}
	if offset := strings.TrimSpace(r.URL.Query().Get("offset")); offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil || parsedOffset < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		input.Offset = &parsedOffset
	}

	out, err := h.listAssessmentsService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]GetAssessmentResponse, 0, len(out.Assessments))
	for _, assessment := range out.Assessments {
		responses = append(responses, GetAssessmentResponse{
			ID:              assessment.ID,
			Title:           assessment.Title,
			Description:     assessment.Description,
			DurationSeconds: assessment.DurationSeconds,
			CreatedBy:       assessment.CreatedBy,
			CreatedAt:       assessment.CreatedAt.Format(http.TimeFormat),
			UpdatedAt:       assessment.UpdatedAt.Format(http.TimeFormat),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ListAssessmentsResponse{Assessments: responses})
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
