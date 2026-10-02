package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/auth/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/auth/internal/adapter/out/userhttp"
	"github.com/aakashloyar/elevate/auth/token"
	"github.com/google/uuid"
)

var ErrInvalidOTP = errors.New("invalid or expired verification code")

type VerifyOTPResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	UserID      string
}

type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, html string) error
}

type Service struct {
	repo          *postgres.Repository
	users         *userhttp.Client
	notifications EmailSender
	secret        string
	ttl           time.Duration
}

func New(repo *postgres.Repository, users *userhttp.Client, notifications EmailSender, secret string, ttl time.Duration) *Service {
	return &Service{repo: repo, users: users, notifications: notifications, secret: secret, ttl: ttl}
}

func (s *Service) TokenExpiresInSeconds() int64 {
	return int64(s.ttl / time.Second)
}

func (s *Service) Register(ctx context.Context, username, email string) (string, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if username == "" || email == "" {
		return "", errors.New("username and email are required")
	}

	userID, err := s.users.Create(ctx, username, email)
	if err != nil {
		return "", err
	}
	if err := s.repo.Save(postgres.Credential{UserID: userID, Username: username, Email: email}); err != nil {
		return "", err
	}
	if err := s.sendOTP(ctx, userID, email); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Service) RequestOTP(ctx context.Context, identifier string) (string, error) {
	credential, err := s.repo.FindByIdentifier(strings.TrimSpace(identifier))
	if err != nil {
		return "", errors.New("account not found")
	}
	if err := s.sendOTP(ctx, credential.UserID, credential.Email); err != nil {
		return "", err
	}
	return credential.UserID, nil
}

func (s *Service) VerifyOTP(userID, code string) (VerifyOTPResult, error) {
	credential, err := s.repo.FindByUserID(userID)
	if err != nil {
		return VerifyOTPResult{}, ErrInvalidOTP
	}

	valid, err := s.repo.ConsumeOTP(credential.UserID, hashCode(code))
	if err != nil {
		return VerifyOTPResult{}, err
	}
	if !valid {
		return VerifyOTPResult{}, ErrInvalidOTP
	}

	now := time.Now().UTC()
	accessToken, err := token.Sign(s.secret, token.Claims{
		Subject:   credential.UserID,
		Email:     credential.Email,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(s.ttl).Unix(),
	})
	if err != nil {
		return VerifyOTPResult{}, err
	}

	return VerifyOTPResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   s.TokenExpiresInSeconds(),
		UserID:      credential.UserID,
	}, nil
}

func (s *Service) sendOTP(ctx context.Context, userID, email string) error {
	var raw [4]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return err
	}

	code := fmt.Sprintf("%06d", int(raw[0])<<8|int(raw[1]))[:6]
	if err := s.repo.SaveOTP(postgres.OTP{
		ID:        uuid.NewString(),
		UserID:    userID,
		CodeHash:  hashCode(code),
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}); err != nil {
		return err
	}

	body := fmt.Sprintf("<p>Your Elevate verification code is <strong>%s</strong>.</p><p>This code expires in 10 minutes.</p>", code)
	return s.notifications.SendEmail(ctx, email, "Your Elevate verification code", body)
}

func hashCode(code string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(code)))
	return hex.EncodeToString(sum[:])
}
