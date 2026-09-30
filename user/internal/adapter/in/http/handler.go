package http

import (
	"encoding/json"
	"net/http"
	"strings"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in"
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

type DeleteUserResponse struct{}

type GetUsersBatchRequest struct {
	UserIDs []string `json:"user_ids"`
}
type UserSummaryResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
type GetUsersBatchResponse struct {
	Users []UserSummaryResponse `json:"users"`
}

type Handler struct {
	createUserService    in.CreateUserService
	getUserService       in.GetUserService
	deleteUserService    in.DeleteUserService
	getUsersBatchService in.GetUsersBatchService
}

func NewHandler(createUserService in.CreateUserService, getUserService in.GetUserService, deleteUserService in.DeleteUserService, getUsersBatchService in.GetUsersBatchService) *Handler {
	return &Handler{createUserService: createUserService, getUserService: getUserService, deleteUserService: deleteUserService, getUsersBatchService: getUsersBatchService}
}

func (h *Handler) GetUsersBatch(w http.ResponseWriter, r *http.Request) {
	var req GetUsersBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	out, err := h.getUsersBatchService.Execute(r.Context(), in.GetUsersBatchInput{UserIDs: req.UserIDs})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users := make([]UserSummaryResponse, 0, len(out.Users))
	for _, user := range out.Users {
		users = append(users, UserSummaryResponse{ID: user.ID, Username: user.Username})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(GetUsersBatchResponse{Users: users})
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

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, userID string) {
	_, err := h.deleteUserService.Execute(r.Context(), in.DeleteUserInput{UserID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(DeleteUserResponse{})
}

func (h *Handler) IsUserRoute(path string) bool {
	return strings.HasPrefix(path, "/users")
}
