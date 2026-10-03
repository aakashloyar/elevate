package in

import "context"

type ListGenerationJobsInput struct {
	Limit  int
	Offset int
	Search string
}

type ListGenerationJobsOutput struct {
	Jobs []GetGenerationJobOutput
}

type ListGenerationJobsService interface {
	Execute(context.Context, ListGenerationJobsInput) (ListGenerationJobsOutput, error)
}
