package service

import (
	"context"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in/user"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
)

type GetUserService struct {
	userRepo out.UserRepository
}

func NewGetUserService(userRepo out.UserRepository) in.GetUserService {
	return &GetUserService{userRepo: userRepo}
}

func (s *GetUserService) Execute(ctx context.Context, input in.GetUserInput) (in.GetUserOutput, error) {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return in.GetUserOutput{}, err
	}

	return in.GetUserOutput{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
