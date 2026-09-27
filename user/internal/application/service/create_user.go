package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
	"github.com/aakashloyar/elevate/user/internal/domain"
	"github.com/lib/pq"
)

type CreateUserService struct {
	userRepo out.UserRepository
	idGen    out.IDGenerator
	clock    out.Clock
}

func NewCreateUserService(userRepo out.UserRepository, idGen out.IDGenerator, clock out.Clock) in.CreateUserService {
	return &CreateUserService{userRepo: userRepo, idGen: idGen, clock: clock}
}

func (s *CreateUserService) Execute(ctx context.Context, input in.CreateUserInput) (in.CreateUserOutput, error) {
	username := strings.TrimSpace(input.Username)
	if username == "" {
		return in.CreateUserOutput{}, errors.New("username is required")
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		return in.CreateUserOutput{}, errors.New("email is required")
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return in.CreateUserOutput{}, errors.New("invalid email")
	}

	now := s.clock.Now()
	user := domain.User{
		ID:        s.idGen.NewID(),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Save(user); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_lower_unique_idx":
				return in.CreateUserOutput{}, errors.New("username already exists")
			case "users_email_lower_unique_idx":
				return in.CreateUserOutput{}, errors.New("email already exists")
			}
		}
		return in.CreateUserOutput{}, err
	}

	return in.CreateUserOutput{UserID: user.ID}, nil
}
