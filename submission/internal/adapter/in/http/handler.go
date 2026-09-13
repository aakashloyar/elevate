package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type CreateSubmissionRequest struct {
	AssessmentID    string `json:"assessment_id"`
	UserID          string `json:"user_id"`
	DurationSeconds int    `json:"duration_seconds"`
}

type CreateSubmissionResponse struct {
	SubmissionID string `json:"submission_id"`
	CreatedAt    string `json:"created_at"`
}

type SaveAnswerRequest struct {
	ProblemID string   `json:"problem_id"`
	Answer    []string `json:"answer"`
}

type SaveAnswerResponse struct{}

type SaveAnswerBatchRequest struct {
	Answers []SaveAnswerRequest `json:"answers"`
}

type SaveAnswerBatchResponse struct {
	SavedCount int                            `json:"saved_count"`
	Errors     []SaveAnswerBatchErrorResponse `json:"errors"`
}

type SaveAnswerBatchErrorResponse struct {
	ProblemID string `json:"problem_id"`
	Message   string `json:"message"`
}

type GetSubmissionResponse struct {
	ID           string                      `json:"id"`
	AssessmentID string                      `json:"assessment_id"`
	UserID       string                      `json:"user_id"`
	Status       string                      `json:"status"`
	StartedAt    *string                     `json:"started_at,omitempty"`
	ExpiresAt    *string                     `json:"expires_at,omitempty"`
	SubmittedAt  *string                     `json:"submitted_at,omitempty"`
	CreatedAt    string                      `json:"created_at"`
	UpdatedAt    string                      `json:"updated_at"`
	Answers      []SubmissionAnswerResponse  `json:"answers"`
	Problems     []SubmissionProblemResponse `json:"problems"`
}

type GetSubmissionStatusResponse struct {
	SubmissionID string  `json:"submission_id"`
	Status       string  `json:"status"`
	ExpiresAt    *string `json:"expires_at"`
}

type UpdateSubmissionStatusRequest struct {
	Status string `json:"status"`
}

type UpdateSubmissionStatusResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
}

type SubmissionAnswerResponse struct {
	ProblemID string   `json:"problem_id"`
	Answer    []string `json:"answer"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type SubmissionProblemResponse struct {
	ProblemID       string   `json:"problem_id"`
	ProblemType     string   `json:"problem_type"`
	OptionIDs       []string `json:"option_ids"`
	OptionTexts     []string `json:"option_texts"`
	AnswerUpdatedAt *string  `json:"answer_updated_at,omitempty"`
}

type Handler struct {
	createSubmissionService       in.CreateSubmissionService
	startSubmissionService        in.StartSubmissionService
	saveAnswerService             in.SaveAnswerService
	saveAnswerBatchService        in.SaveAnswerBatchService
	getSubmissionService          in.GetSubmissionService
	getSubmissionStatusService    in.GetSubmissionStatusService
	submitSubmissionService       in.SubmitSubmissionService
	updateSubmissionStatusService in.UpdateSubmissionStatusService
}

func NewHandler(createSubmissionService in.CreateSubmissionService, startSubmissionService in.StartSubmissionService, saveAnswerService in.SaveAnswerService, saveAnswerBatchService in.SaveAnswerBatchService, getSubmissionService in.GetSubmissionService, getSubmissionStatusService in.GetSubmissionStatusService, submitSubmissionService in.SubmitSubmissionService, updateSubmissionStatusService in.UpdateSubmissionStatusService) *Handler {
	return &Handler{
		createSubmissionService:       createSubmissionService,
		startSubmissionService:        startSubmissionService,
		saveAnswerService:             saveAnswerService,
		saveAnswerBatchService:        saveAnswerBatchService,
		getSubmissionService:          getSubmissionService,
		getSubmissionStatusService:    getSubmissionStatusService,
		submitSubmissionService:       submitSubmissionService,
		updateSubmissionStatusService: updateSubmissionStatusService,
	}
}

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.createSubmissionService.Execute(r.Context(), in.CreateSubmissionInput{AssessmentID: req.AssessmentID, UserID: req.UserID, DurationSeconds: req.DurationSeconds})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateSubmissionResponse{SubmissionID: out.SubmissionID, CreatedAt: out.CreatedAt})
}

func (h *Handler) SaveAnswer(w http.ResponseWriter, r *http.Request, submissionID string) {
	var req SaveAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.saveAnswerService.Execute(r.Context(), in.SaveAnswerInput{SubmissionID: submissionID, ProblemID: req.ProblemID, Answer: req.Answer}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SaveAnswerBatch(w http.ResponseWriter, r *http.Request, submissionID string) {
	var req SaveAnswerBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	answers := make([]in.SaveAnswerBatchItem, 0, len(req.Answers))
	for _, item := range req.Answers {
		answers = append(answers, in.SaveAnswerBatchItem{ProblemID: item.ProblemID, Answer: item.Answer})
	}

	out, err := h.saveAnswerBatchService.Execute(r.Context(), in.SaveAnswerBatchInput{SubmissionID: submissionID, Answers: answers})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := make([]SaveAnswerBatchErrorResponse, 0, len(out.Errors))
	for _, item := range out.Errors {
		errors = append(errors, SaveAnswerBatchErrorResponse{
			ProblemID: item.ProblemID,
			Message:   item.Message,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SaveAnswerBatchResponse{
		SavedCount: out.SavedCount,
		Errors:     errors,
	})
}

func (h *Handler) GetSubmissionByID(w http.ResponseWriter, r *http.Request, submissionID string) {
	out, err := h.getSubmissionService.Execute(r.Context(), in.GetSubmissionInput{SubmissionID: submissionID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var submittedAt *string
	if out.SubmittedAt != nil {
		value := out.SubmittedAt.Format(http.TimeFormat)
		submittedAt = &value
	}
	var startedAt *string
	if out.StartedAt != nil {
		value := out.StartedAt.Format(http.TimeFormat)
		startedAt = &value
	}
	var expiresAt *string
	if out.ExpiresAt != nil {
		value := out.ExpiresAt.Format(http.TimeFormat)
		expiresAt = &value
	}

	answers := make([]SubmissionAnswerResponse, 0, len(out.Answers))
	for _, ans := range out.Answers {
		answers = append(answers, SubmissionAnswerResponse{
			ProblemID: ans.ProblemID,
			Answer:    ans.Answer,
			CreatedAt: ans.CreatedAt.Format(http.TimeFormat),
			UpdatedAt: ans.UpdatedAt.Format(http.TimeFormat),
		})
	}
	problems := make([]SubmissionProblemResponse, 0, len(out.Problems))
	for _, problem := range out.Problems {
		var answerUpdatedAt *string
		if problem.AnswerUpdatedAt != nil {
			value := problem.AnswerUpdatedAt.Format(http.TimeFormat)
			answerUpdatedAt = &value
		}
		problems = append(problems, SubmissionProblemResponse{
			ProblemID:       problem.ProblemID,
			ProblemType:     string(problem.ProblemType),
			OptionIDs:       problem.OptionIDs,
			OptionTexts:     problem.OptionTexts,
			AnswerUpdatedAt: answerUpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetSubmissionResponse{
		ID:           out.ID,
		AssessmentID: out.AssessmentID,
		UserID:       out.UserID,
		Status:       string(out.Status),
		StartedAt:    startedAt,
		ExpiresAt:    expiresAt,
		SubmittedAt:  submittedAt,
		CreatedAt:    out.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:    out.UpdatedAt.Format(http.TimeFormat),
		Answers:      answers,
		Problems:     problems,
	})
}

func (h *Handler) GetSubmissionStatus(w http.ResponseWriter, r *http.Request, submissionID string) {
	out, err := h.getSubmissionStatusService.Execute(r.Context(), in.GetSubmissionStatusInput{SubmissionID: submissionID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var expiresAt *string
	if out.ExpiresAt != nil {
		value := out.ExpiresAt.Format(http.TimeFormat)
		expiresAt = &value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetSubmissionStatusResponse{
		SubmissionID: out.SubmissionID,
		Status:       string(out.Status),
		ExpiresAt:    expiresAt,
	})
}

func (h *Handler) UpdateSubmissionStatus(w http.ResponseWriter, r *http.Request, submissionID string) {
	var req UpdateSubmissionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.updateSubmissionStatusService.Execute(r.Context(), in.UpdateSubmissionStatusInput{
		SubmissionID: submissionID,
		Status:       domain.SubmissionStatus(req.Status),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UpdateSubmissionStatusResponse{
		SubmissionID: out.SubmissionID,
		Status:       string(out.Status),
	})
}

func (h *Handler) SubmitSubmission(w http.ResponseWriter, r *http.Request, submissionID string) {
	if err := h.submitSubmissionService.Execute(r.Context(), in.SubmitSubmissionInput{SubmissionID: submissionID}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) StartSubmission(w http.ResponseWriter, r *http.Request, submissionID string) {
	if err := h.startSubmissionService.Execute(r.Context(), in.StartSubmissionInput{SubmissionID: submissionID}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) IsSubmissionRoute(path string) bool {
	return strings.HasPrefix(path, "/submissions")
}
