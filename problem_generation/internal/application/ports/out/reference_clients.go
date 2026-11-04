package out

import "context"

type UserClient interface {
	Exists(ctx context.Context, userID string) error
}

type AssessmentClient interface {
	Exists(ctx context.Context, assessmentID string) error
	GetTitles(ctx context.Context, assessmentIDs []string) (map[string]string, error)
}

type DocumentSummary struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by"`
	Status    string `json:"status"`
}

type DocumentClient interface {
	GetDocument(ctx context.Context, documentID string) (DocumentSummary, error)
	GetContent(ctx context.Context, documentID string) (string, error)
}
