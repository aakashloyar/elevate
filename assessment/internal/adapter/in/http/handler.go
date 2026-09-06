package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	outports "github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	assessmentsvc "github.com/aakashloyar/elevate/assessment/internal/application/service"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
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
	createAssessmentService      in.CreateAssessmentService
	listAssessmentsService       in.ListAssessmentsService
	getAssessmentService         in.GetAssessmentService
	deleteAssessmentService      in.DeleteAssessmentService
	createProblemService         in.CreateAssessmentProblemService
	getAssessmentProblemsService in.GetAssessmentProblemsService
	getMarkingSchemeService      in.GetAssessmentMarkingSchemeService
	upsertMarkingSchemeService   in.UpsertAssessmentMarkingSchemeService
	createMarkingSchemeService   in.CreateAssessmentMarkingSchemeService
}

func NewHandler(createAssessmentService in.CreateAssessmentService, listAssessmentsService in.ListAssessmentsService, getAssessmentService in.GetAssessmentService, deleteAssessmentService in.DeleteAssessmentService, createProblemService in.CreateAssessmentProblemService, getAssessmentProblemsService in.GetAssessmentProblemsService, getMarkingSchemeService in.GetAssessmentMarkingSchemeService, upsertMarkingSchemeService in.UpsertAssessmentMarkingSchemeService, createMarkingSchemeService in.CreateAssessmentMarkingSchemeService) *Handler {
	return &Handler{
		createAssessmentService:      createAssessmentService,
		listAssessmentsService:       listAssessmentsService,
		getAssessmentService:         getAssessmentService,
		deleteAssessmentService:      deleteAssessmentService,
		createProblemService:         createProblemService,
		getAssessmentProblemsService: getAssessmentProblemsService,
		getMarkingSchemeService:      getMarkingSchemeService,
		upsertMarkingSchemeService:   upsertMarkingSchemeService,
		createMarkingSchemeService:   createMarkingSchemeService,
	}
}

type ListAssessmentsResponse struct {
	Assessments []GetAssessmentResponse `json:"assessments"`
}

type GetAssessmentProblemsResponse struct {
	ProblemIDs []string `json:"problem_ids"`
}

type CreateAssessmentProblemRequest struct {
	CreatedBy  string                               `json:"created_by"`
	Title      string                               `json:"title"`
	Statement  string                               `json:"statement"`
	Type       string                               `json:"type"`
	Difficulty string                               `json:"difficulty"`
	SourceType string                               `json:"source_type"`
	Options    []CreateAssessmentProblemOptionInput `json:"options"`
	Tags       []string                             `json:"tags"`
}

type CreateAssessmentProblemOptionInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateAssessmentProblemResponse struct {
	ProblemID string `json:"problem_id"`
}

type AssessmentMarkingSchemeRequest struct {
	Single    MarksRequest `json:"single"`
	Multiple  MarksRequest `json:"multiple"`
	Numerical MarksRequest `json:"numerical"`
}

type MarksRequest struct {
	Correct   float64 `json:"correct"`
	Incorrect float64 `json:"incorrect"`
	Skipped   float64 `json:"skipped"`
}

type AssessmentMarkingSchemeResponse struct {
	AssessmentID string `json:"assessment_id"`
	AssessmentMarkingSchemeRequest
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

func (h *Handler) GetAssessmentMarkingScheme(w http.ResponseWriter, r *http.Request, assessmentID string) {
	markingScheme, err := h.getMarkingSchemeService.Execute(r.Context(), in.GetAssessmentMarkingSchemeInput{AssessmentID: assessmentID})
	if err != nil {
		http.Error(w, "assessment marking scheme not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toAssessmentMarkingSchemeResponse(markingScheme))
}

func (h *Handler) GetAssessmentProblems(w http.ResponseWriter, r *http.Request, assessmentID string) {
	out, err := h.getAssessmentProblemsService.Execute(r.Context(), in.GetAssessmentProblemsInput{AssessmentID: assessmentID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GetAssessmentProblemsResponse{ProblemIDs: out.ProblemIDs})
}

func (h *Handler) PutAssessmentMarkingScheme(w http.ResponseWriter, r *http.Request, assessmentID string) {
	var req AssessmentMarkingSchemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.upsertMarkingSchemeService.Execute(r.Context(), toUpsertAssessmentMarkingSchemeInput(assessmentID, req))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "assessment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toAssessmentMarkingSchemeResponseFromRequest(assessmentID, req))
}

func (h *Handler) CreateAssessmentMarkingScheme(w http.ResponseWriter, r *http.Request, assessmentID string) {
	var req AssessmentMarkingSchemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.createMarkingSchemeService.Execute(r.Context(), toUpsertAssessmentMarkingSchemeInput(assessmentID, req))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "assessment not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, assessmentsvc.ErrAssessmentMarkingSchemeAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toAssessmentMarkingSchemeResponseFromRequest(assessmentID, req))
}

func toUpsertAssessmentMarkingSchemeInput(assessmentID string, req AssessmentMarkingSchemeRequest) in.UpsertAssessmentMarkingSchemeInput {
	return in.UpsertAssessmentMarkingSchemeInput{
		AssessmentID: assessmentID,
		Single:       toMarksInput(req.Single),
		Multiple:     toMarksInput(req.Multiple),
		Numerical:    toMarksInput(req.Numerical),
	}
}

func toMarksInput(req MarksRequest) in.MarksInput {
	return in.MarksInput{
		Correct:   req.Correct,
		Incorrect: req.Incorrect,
		Skipped:   req.Skipped,
	}
}

func toAssessmentMarkingSchemeResponse(markingScheme domain.AssessmentMarkingScheme) AssessmentMarkingSchemeResponse {
	return AssessmentMarkingSchemeResponse{
		AssessmentID: markingScheme.AssessmentID,
		AssessmentMarkingSchemeRequest: AssessmentMarkingSchemeRequest{
			Single:    toMarksRequest(markingScheme.Single),
			Multiple:  toMarksRequest(markingScheme.Multiple),
			Numerical: toMarksRequest(markingScheme.Numerical),
		},
	}
}

func toAssessmentMarkingSchemeResponseFromRequest(assessmentID string, req AssessmentMarkingSchemeRequest) AssessmentMarkingSchemeResponse {
	return AssessmentMarkingSchemeResponse{AssessmentID: assessmentID, AssessmentMarkingSchemeRequest: req}
}

func toMarksRequest(marks domain.Marks) MarksRequest {
	return MarksRequest{
		Correct:   marks.Correct,
		Incorrect: marks.Incorrect,
		Skipped:   marks.Skipped,
	}
}

func (h *Handler) CreateAssessmentProblem(w http.ResponseWriter, r *http.Request, assessmentID string) {
	var req CreateAssessmentProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	options := make([]in.CreateAssessmentProblemOptionInput, 0, len(req.Options))
	for _, option := range req.Options {
		options = append(options, in.CreateAssessmentProblemOptionInput{Text: option.Text, IsCorrect: option.IsCorrect})
	}

	out, err := h.createProblemService.Execute(r.Context(), in.CreateAssessmentProblemInput{
		AssessmentID: assessmentID,
		CreatedBy:    req.CreatedBy,
		Title:        req.Title,
		Statement:    req.Statement,
		Type:         req.Type,
		Difficulty:   req.Difficulty,
		SourceType:   req.SourceType,
		Options:      options,
		Tags:         req.Tags,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "assessment not found", http.StatusNotFound)
			return
		}

		var problemErr *outports.ProblemClientError
		if errors.As(err, &problemErr) {
			if problemErr.StatusCode >= 400 && problemErr.StatusCode < 500 {
				http.Error(w, problemErr.Message, http.StatusBadRequest)
				return
			}
			http.Error(w, "problem service unavailable", http.StatusBadGateway)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateAssessmentProblemResponse{ProblemID: out.ProblemID})
}

func (h *Handler) IsAssessmentRoute(path string) bool {
	return strings.HasPrefix(path, "/assessments")
}
