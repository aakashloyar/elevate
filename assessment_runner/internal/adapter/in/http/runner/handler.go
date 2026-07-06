package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in/runner"
)

type StartSessionRequest struct {
	AssessmentID string `json:"assessment_id"`
	UserID       string `json:"user_id"`
}

type StartSessionResponse struct {
	SessionID      string `json:"session_id"`
	SubmissionID   string `json:"submission_id"`
	RemainingTime  int    `json:"remaining_time"`
	TotalQuestions int    `json:"total_questions"`
}

type GetSessionResponse struct {
	SessionID      string `json:"session_id"`
	RemainingTime  int    `json:"remaining_time"`
	TotalQuestions int    `json:"total_questions"`
	Status         string `json:"status"`
}

type GetQuestionsResponse struct {
	Questions []in.QuestionView `json:"questions"`
}

type Handler struct {
	startSessionService  in.StartSessionService
	getSessionService    in.GetSessionService
	getQuestionsService  in.GetQuestionsService
	submitSessionService in.SubmitSessionService
}

func NewHandler(startSessionService in.StartSessionService) *Handler {
	return &Handler{startSessionService: startSessionService}
}

func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.startSessionService.Execute(r.Context(), in.StartSessionInput{AssessmentID: req.AssessmentID, UserID: req.UserID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(StartSessionResponse{SessionID: out.SessionID, SubmissionID: out.SubmissionID, RemainingTime: out.RemainingTime, TotalQuestions: out.TotalQuestions})
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	out, err := h.getSessionService.Execute(r.Context(), in.GetSessionInput{SessionID: sessionID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetSessionResponse{SessionID: out.SessionID, RemainingTime: out.RemainingTime, TotalQuestions: out.TotalQuestions, Status: out.Status})
}

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request, sessionID string) {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10
	}

	out, err := h.getQuestionsService.Execute(r.Context(), in.GetQuestionsInput{SessionID: sessionID, Offset: offset, Limit: limit})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetQuestionsResponse{Questions: out.Questions})
}

func (h *Handler) SubmitSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	if err := h.submitSessionService.Execute(r.Context(), in.SubmitSessionInput{SessionID: sessionID}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) IsRunnerRoute(path string) bool {
	return strings.HasPrefix(path, "/sessions")
}
