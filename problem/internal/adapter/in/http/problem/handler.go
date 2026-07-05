package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
)

type CreateProblemRequest struct {
	CreatedBy  string                     `json:"created_by"`
	Title      string                     `json:"title"`
	Statement  string                     `json:"statement"`
	Type       string                     `json:"type"`
	Difficulty string                     `json:"difficulty"`
	SourceType string                     `json:"source_type"`
	Status     string                     `json:"status"`
	Options    []CreateProblemOptionInput `json:"options"`
	Tags       []string                   `json:"tags"`
}

type CreateProblemOptionInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateProblemResponse struct {
	ProblemID string `json:"problem_id"`
}

type GetProblemResponse struct {
	ID         string                  `json:"id"`
	CreatedBy  string                  `json:"created_by"`
	Title      string                  `json:"title"`
	Statement  string                  `json:"statement"`
	Type       string                  `json:"type"`
	Difficulty string                  `json:"difficulty"`
	SourceType string                  `json:"source_type"`
	Status     string                  `json:"status"`
	Options    []ProblemOptionResponse `json:"options"`
	Tags       []string                `json:"tags"`
	CreatedAt  string                  `json:"created_at"`
	UpdatedAt  string                  `json:"updated_at"`
}

type ProblemOptionResponse struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type ListProblemsResponse struct {
	Problems []ListProblemItemResponse `json:"problems"`
}

type ListProblemItemResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	Difficulty string `json:"difficulty"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type Handler struct {
	createProblemService in.CreateProblemService
	getProblemService    in.GetProblemService
	listProblemsService  in.ListProblemsService
	updateProblemService in.UpdateProblemService
	deleteProblemService in.DeleteProblemService
}

func NewHandler(createProblemService in.CreateProblemService, getProblemService in.GetProblemService, listProblemsService in.ListProblemsService, updateProblemService in.UpdateProblemService, deleteProblemService in.DeleteProblemService) *Handler {
	return &Handler{
		createProblemService: createProblemService,
		getProblemService:    getProblemService,
		listProblemsService:  listProblemsService,
		updateProblemService: updateProblemService,
		deleteProblemService: deleteProblemService,
	}
}

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var req CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	options := make([]in.CreateProblemOptionInput, 0, len(req.Options))
	for _, opt := range req.Options {
		options = append(options, in.CreateProblemOptionInput{Text: opt.Text, IsCorrect: opt.IsCorrect})
	}

	out, err := h.createProblemService.Execute(r.Context(), in.CreateProblemInput{
		CreatedBy:  req.CreatedBy,
		Title:      req.Title,
		Statement:  req.Statement,
		Type:       req.Type,
		Difficulty: req.Difficulty,
		SourceType: req.SourceType,
		Status:     req.Status,
		Options:    options,
		Tags:       req.Tags,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateProblemResponse{ProblemID: out.ProblemID})
}

func (h *Handler) GetProblemByID(w http.ResponseWriter, r *http.Request, problemID string) {
	out, err := h.getProblemService.Execute(r.Context(), in.GetProblemInput{ProblemID: problemID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	options := make([]ProblemOptionResponse, 0, len(out.Options))
	for _, option := range out.Options {
		options = append(options, ProblemOptionResponse{ID: option.ID, Text: option.Text, IsCorrect: option.IsCorrect})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetProblemResponse{
		ID:         out.ID,
		CreatedBy:  out.CreatedBy,
		Title:      out.Title,
		Statement:  out.Statement,
		Type:       out.Type,
		Difficulty: out.Difficulty,
		SourceType: out.SourceType,
		Status:     out.Status,
		Options:    options,
		Tags:       out.Tags,
		CreatedAt:  out.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:  out.UpdatedAt.Format(http.TimeFormat),
	})
}

func (h *Handler) ListProblems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}

	out, err := h.listProblemsService.Execute(r.Context(), in.ListProblemsInput{
		Offset: offset,
		Limit:  limit,
		Type:   query.Get("type"),
		Status: query.Get("status"),
		Tag:    query.Get("tag"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items := make([]ListProblemItemResponse, 0, len(out.Problems))
	for _, item := range out.Problems {
		items = append(items, ListProblemItemResponse{ID: item.ID, Title: item.Title, Type: item.Type, Difficulty: item.Difficulty, Status: item.Status, CreatedAt: item.CreatedAt})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ListProblemsResponse{Problems: items})
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request, problemID string) {
	var req CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	options := make([]in.CreateProblemOptionInput, 0, len(req.Options))
	for _, opt := range req.Options {
		options = append(options, in.CreateProblemOptionInput{Text: opt.Text, IsCorrect: opt.IsCorrect})
	}

	if err := h.updateProblemService.Execute(r.Context(), in.UpdateProblemInput{
		ProblemID:  problemID,
		CreatedBy:  req.CreatedBy,
		Title:      req.Title,
		Statement:  req.Statement,
		Type:       req.Type,
		Difficulty: req.Difficulty,
		SourceType: req.SourceType,
		Status:     req.Status,
		Options:    options,
		Tags:       req.Tags,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteProblem(w http.ResponseWriter, r *http.Request, problemID string) {
	if err := h.deleteProblemService.Execute(r.Context(), in.DeleteProblemInput{ProblemID: problemID}); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) IsProblemRoute(path string) bool {
	return strings.HasPrefix(path, "/problems")
}
