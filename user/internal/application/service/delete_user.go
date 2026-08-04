package service

import (
	"context"

	"github.com/aakashloyar/elevate/user/config"
	in "github.com/aakashloyar/elevate/user/internal/application/ports/in"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
)

type DeleteUserService struct {
	userRepo out.UserRepository
}

func NewDeleteUserService(userRepo out.UserRepository) in.DeleteUserService {
	return &DeleteUserService{userRepo: userRepo}
}

func (s *DeleteUserService) Execute(ctx context.Context, input in.DeleteUserInput) (in.DeleteUserOutput, error) {
	result, err := s.userRepo.Delete(input.UserID)
	if err != nil {
		return in.DeleteUserOutput{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return in.DeleteUserOutput{}, err
	}

	if rowsAffected == 0 {
		return in.DeleteUserOutput{}, config.ErrUserNotFound
	}
	return in.DeleteUserOutput{}, nil
}
