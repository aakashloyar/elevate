package service

import (
	"context"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in/user"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
)

type DeleteUserService struct {
	userRepo out.UserRepository
}

func NewDeleteUserService(userRepo out.UserRepository) in.DeleteUserService {
	return &DeleteUserService{userRepo: userRepo}
}

func (s *DeleteUserService) Execute(ctx context.Context, input in.DeleteUserInput) (in.DeleteUserOutput, error) {
	if err := s.userRepo.Delete(input.UserID); err != nil {
		return in.DeleteUserOutput{}, err
	}

	return in.DeleteUserOutput{Deleted: true}, nil
}
