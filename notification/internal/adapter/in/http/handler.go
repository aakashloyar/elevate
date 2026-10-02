package http

import (
	"encoding/json"
	"net/http"

	"github.com/aakashloyar/elevate/notification/internal/email"
)

type Handler struct {
	sender *email.Sender
}

func NewHandler(sender *email.Sender) *Handler { return &Handler{sender: sender} }

type emailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req emailRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.sender.Send(req.To, req.Subject, req.Text, req.HTML); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
