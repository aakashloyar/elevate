package service

import (
	"context"

	in "github.com/aakashloyar/elevate/user/internal/application/ports/in"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
)

type GetUsersBatchService struct{ repository out.UserRepository }

func NewGetUsersBatchService(repository out.UserRepository) in.GetUsersBatchService {
	return &GetUsersBatchService{repository: repository}
}

func (s *GetUsersBatchService) Execute(_ context.Context, input in.GetUsersBatchInput) (in.GetUsersBatchOutput, error) {
	users, err := s.repository.FindByIDs(input.UserIDs)
	if err != nil {
		return in.GetUsersBatchOutput{}, err
	}
	outUsers := make([]in.UserSummary, 0, len(users))
	for _, user := range users {
		outUsers = append(outUsers, in.UserSummary{ID: user.ID, Username: user.Username})
	}
	return in.GetUsersBatchOutput{Users: outUsers}, nil
}
