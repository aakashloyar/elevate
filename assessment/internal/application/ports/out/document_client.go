package out

import "context"

type DocumentSummary struct {
	ID        string `json:"id"`
	CreatedBy string `json:"created_by"`
	Status    string `json:"status"`
}

type DocumentClient interface {
	GetDocuments(ctx context.Context, documentIDs []string) ([]DocumentSummary, error)
}
