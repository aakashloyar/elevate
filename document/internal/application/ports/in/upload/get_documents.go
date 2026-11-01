package upload

import (
	"context"
)

type GetDocumentsInput struct {
	DocumentIDs []string
}

type GetDocumentsOutput struct {
	Documents []DocumentSummary
}

type GetDocumentsService interface {
	Execute(context.Context, GetDocumentsInput) (GetDocumentsOutput, error)
}
