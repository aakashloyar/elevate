package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in/runner"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/out"
)

type RunnerService struct {
	clientFactory out.ClientFactory
}

func NewRunnerService(clientFactory out.ClientFactory) in.StartSessionService {
	return &RunnerService{clientFactory: clientFactory}
}

func (s *RunnerService) Execute(ctx context.Context, input in.StartSessionInput) (in.StartSessionOutput, error) {
	assessmentClient := s.clientFactory.AssessmentClient()
	problemClient := s.clientFactory.ProblemClient()
	submissionClient := s.clientFactory.SubmissionClient()

	assessmentResp, err := assessmentClient.Get(ctx, "/assessments/"+input.AssessmentID)
	if err != nil {
		return in.StartSessionOutput{}, err
	}
	defer assessmentResp.Body.Close()
	if assessmentResp.StatusCode != http.StatusOK {
		return in.StartSessionOutput{}, fmt.Errorf("assessment lookup failed")
	}

	submissionReq := map[string]string{"assessment_id": input.AssessmentID, "user_id": input.UserID}
	submissionResp, err := submissionClient.Post(ctx, "/submissions", submissionReq)
	if err != nil {
		return in.StartSessionOutput{}, err
	}
	defer submissionResp.Body.Close()
	if submissionResp.StatusCode != http.StatusCreated {
		return in.StartSessionOutput{}, fmt.Errorf("submission creation failed")
	}

	var submissionBody struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := json.NewDecoder(submissionResp.Body).Decode(&submissionBody); err != nil {
		return in.StartSessionOutput{}, err
	}

	var assessmentBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(assessmentResp.Body).Decode(&assessmentBody); err != nil {
		return in.StartSessionOutput{}, err
	}

	questionResp, err := problemClient.Get(ctx, "/problems?limit=100")
	if err != nil {
		return in.StartSessionOutput{}, err
	}
	defer questionResp.Body.Close()
	if questionResp.StatusCode != http.StatusOK {
		return in.StartSessionOutput{}, fmt.Errorf("problem lookup failed")
	}

	var problemsBody struct {
		Problems []struct {
			ID string `json:"id"`
		} `json:"problems"`
	}
	if err := json.NewDecoder(questionResp.Body).Decode(&problemsBody); err != nil {
		return in.StartSessionOutput{}, err
	}

	return in.StartSessionOutput{
		SessionID:      fmt.Sprintf("session-%s", input.AssessmentID),
		SubmissionID:   submissionBody.SubmissionID,
		RemainingTime:  7200,
		TotalQuestions: len(problemsBody.Problems),
	}, nil
}

func (s *RunnerService) GetSession(ctx context.Context, input in.GetSessionInput) (in.GetSessionOutput, error) {
	return in.GetSessionOutput{SessionID: input.SessionID, RemainingTime: 7200, TotalQuestions: 100, Status: "ACTIVE"}, nil
}

func (s *RunnerService) GetQuestions(ctx context.Context, input in.GetQuestionsInput) (in.GetQuestionsOutput, error) {
	problemClient := s.clientFactory.ProblemClient()
	resp, err := problemClient.Get(ctx, "/problems?offset="+strconv.Itoa(input.Offset)+"&limit="+strconv.Itoa(input.Limit))
	if err != nil {
		return in.GetQuestionsOutput{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return in.GetQuestionsOutput{}, fmt.Errorf("problem lookup failed")
	}

	var body struct {
		Problems []struct {
			ID         string `json:"id"`
			Title      string `json:"title"`
			Type       string `json:"type"`
			Difficulty string `json:"difficulty"`
			Status     string `json:"status"`
		} `json:"problems"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return in.GetQuestionsOutput{}, err
	}

	questions := make([]in.QuestionView, 0, len(body.Problems))
	for _, problem := range body.Problems {
		questions = append(questions, in.QuestionView{ID: problem.ID, Title: problem.Title, Type: problem.Type, Difficulty: problem.Difficulty})
	}
	return in.GetQuestionsOutput{Questions: questions}, nil
}

func (s *RunnerService) SubmitSession(ctx context.Context, input in.SubmitSessionInput) error {
	return nil
}
