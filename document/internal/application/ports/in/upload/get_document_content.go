package upload

import "context"

type GetDocumentContentInput struct {
	DocumentID string
}

type GetDocumentContentOutput struct {
	DocumentID string
	Content    string
}

type GetDocumentContentService interface {
	Execute(context.Context, GetDocumentContentInput) (GetDocumentContentOutput, error)
}
