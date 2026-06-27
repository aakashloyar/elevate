package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in/user"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type GetUserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type Handler struct {
	createUserService in.CreateUserService
	getUserService    in.GetUserService
}

func NewHandler(createUserService in.CreateUserService, getUserService in.GetUserService) *Handler {
	return &Handler{createUserService: createUserService, getUserService: getUserService}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.createUserService.Execute(r.Context(), in.CreateUserInput{Username: req.Username, Email: req.Email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateUserResponse{UserID: out.UserID})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, userID string) {
	out, err := h.getUserService.Execute(r.Context(), in.GetUserInput{UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetUserResponse{
		ID:        out.ID,
		Username:  out.Username,
		Email:     out.Email,
		CreatedAt: out.CreatedAt.Format(http.TimeFormat),
	})
}

func (h *Handler) IsUserRoute(path string) bool {
	return strings.HasPrefix(path, "/users")
}
