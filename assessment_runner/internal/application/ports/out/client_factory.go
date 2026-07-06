package out

import (
	"context"
	"net/http"
)

type ClientFactory interface {
	AssessmentClient() ServiceClient
	ProblemClient() ServiceClient
	SubmissionClient() ServiceClient
}

type ServiceClient interface {
	Get(ctx context.Context, path string) (*http.Response, error)
	Post(ctx context.Context, path string, body any) (*http.Response, error)
}
