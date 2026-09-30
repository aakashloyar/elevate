package submission

import (
	"context"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type ListSubmissionsService struct {
	repository       out.SubmissionRepository
	userClient       out.UserClient
	assessmentClient out.AssessmentClient
}

func NewListSubmissionsService(repository out.SubmissionRepository, userClient out.UserClient, assessmentClient out.AssessmentClient) in.ListSubmissionsService {
	return &ListSubmissionsService{repository: repository, userClient: userClient, assessmentClient: assessmentClient}
}

func (s *ListSubmissionsService) Execute(ctx context.Context, input in.ListSubmissionsInput) (in.ListSubmissionsOutput, error) {
	var submissions []domain.Submission
	var err error
	if strings.TrimSpace(input.UserID) != "" {
		submissions, err = s.repository.ListByUserID(strings.TrimSpace(input.UserID), input.Limit, input.Offset)
	} else {
		submissions, err = s.repository.ListAll(input.Limit, input.Offset)
	}
	if err != nil {
		return in.ListSubmissionsOutput{}, err
	}
	userIDs := make([]string, 0, len(submissions))
	for _, submission := range submissions {
		userIDs = append(userIDs, submission.UserID)
	}
	usernames, err := s.userClient.GetUsernames(ctx, userIDs)
	if err != nil {
		return in.ListSubmissionsOutput{}, err
	}
	assessmentIDs := make([]string, 0, len(submissions))
	for _, submission := range submissions {
		assessmentIDs = append(assessmentIDs, submission.AssessmentID)
	}
	assessmentTitles, err := s.assessmentClient.GetAssessmentTitles(ctx, assessmentIDs)
	if err != nil {
		return in.ListSubmissionsOutput{}, err
	}
	items := make([]in.ListSubmissionItem, 0, len(submissions))
	for _, submission := range submissions {
		items = append(items, in.ListSubmissionItem{Submission: submission, UserName: usernames[submission.UserID], AssessmentTitle: assessmentTitles[submission.AssessmentID]})
	}
	return in.ListSubmissionsOutput{Submissions: items}, nil
}
