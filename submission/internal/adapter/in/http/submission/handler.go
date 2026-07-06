package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in/submission"
)

type CreateSubmissionRequest struct {
	AssessmentID string `json:"assessment_id"`
	UserID       string `json:"user_id"`
}

type CreateSubmissionResponse struct {
	SubmissionID string `json:"submission_id"`
	StartedAt    string `json:"started_at"`
}

type SubmitAnswerRequest struct {
	ProblemID string   `json:"problem_id"`
	Answer    []string `json:"answer"`
}

type SubmitAnswerResponse struct{}

type GetSubmissionResponse struct {
	ID           string                     `json:"id"`
	AssessmentID string                     `json:"assessment_id"`
	UserID       string                     `json:"user_id"`
	Status       string                     `json:"status"`
	StartedAt    string                     `json:"started_at"`
	SubmittedAt  *string                    `json:"submitted_at,omitempty"`
	Answers      []SubmissionAnswerResponse `json:"answers"`
}

type SubmissionAnswerResponse struct {
	ProblemID  string   `json:"problem_id"`
	Answer     []string `json:"answer"`
	AnsweredAt string   `json:"answered_at"`
}

type Handler struct {
	createSubmissionService in.CreateSubmissionService
	submitAnswerService     in.SubmitAnswerService
	getSubmissionService    in.GetSubmissionService
}

func NewHandler(createSubmissionService in.CreateSubmissionService, submitAnswerService in.SubmitAnswerService, getSubmissionService in.GetSubmissionService) *Handler {
	return &Handler{
		createSubmissionService: createSubmissionService,
		submitAnswerService:     submitAnswerService,
		getSubmissionService:    getSubmissionService,
	}
}

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.createSubmissionService.Execute(r.Context(), in.CreateSubmissionInput{AssessmentID: req.AssessmentID, UserID: req.UserID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateSubmissionResponse{SubmissionID: out.SubmissionID, StartedAt: out.StartedAt})
}

func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request, submissionID string) {
	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.submitAnswerService.Execute(r.Context(), in.SubmitAnswerInput{SubmissionID: submissionID, ProblemID: req.ProblemID, Answer: req.Answer}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	answers := make([]SubmissionAnswerResponse, 0, len(out.Answers))
	for _, ans := range out.Answers {
		answers = append(answers, SubmissionAnswerResponse{ProblemID: ans.ProblemID, Answer: ans.Answer, AnsweredAt: ans.AnsweredAt.Format(http.TimeFormat)})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetSubmissionResponse{
		ID:           out.ID,
		AssessmentID: out.AssessmentID,
		UserID:       out.UserID,
		Status:       out.Status,
		StartedAt:    out.StartedAt.Format(http.TimeFormat),
		SubmittedAt:  submittedAt,
		Answers:      answers,
	})
}

func (h *Handler) IsSubmissionRoute(path string) bool {
	return strings.HasPrefix(path, "/submissions")
}
