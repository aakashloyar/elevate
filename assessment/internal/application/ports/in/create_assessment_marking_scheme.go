package assessment

import "context"

type CreateAssessmentMarkingSchemeService interface {
	Execute(ctx context.Context, input UpsertAssessmentMarkingSchemeInput) error
}
