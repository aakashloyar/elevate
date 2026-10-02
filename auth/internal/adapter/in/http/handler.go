package http

import (
	"encoding/json"
	"net/http"

	"github.com/aakashloyar/elevate/auth/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	UserID               string `json:"user_id"`
	VerificationRequired bool   `json:"verification_required"`
}

type RequestOTPRequest struct {
	Identifier string `json:"identifier"`
}

type RequestOTPResponse struct {
	UserID               string `json:"user_id"`
	VerificationRequired bool   `json:"verification_required"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}

type VerifyOTPResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.service.Register(r.Context(), req.Username, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(RegisterResponse{UserID: id, VerificationRequired: true})
}
func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req RequestOTPRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := h.service.RequestOTP(r.Context(), req.Identifier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(RequestOTPResponse{UserID: id, VerificationRequired: true})
}
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req VerifyOTPRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result, err := h.service.VerifyOTP(req.UserID, req.OTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(VerifyOTPResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
		UserID:      result.UserID,
	})
}
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
