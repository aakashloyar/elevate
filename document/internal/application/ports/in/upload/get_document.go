package upload

import "context"

type GetDocumentInput struct {
	DocumentID string
}

type GetDocumentOutput struct {
	Document DocumentSummary
}

type GetDocumentService interface {
	Execute(context.Context, GetDocumentInput) (GetDocumentOutput, error)
}
