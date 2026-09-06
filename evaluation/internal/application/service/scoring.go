package evaluation

import "github.com/aakashloyar/elevate/evaluation/internal/domain"

func EvaluateAnswer(problem domain.Problem, selectedOptions []string, scheme domain.MarkingScheme) domain.QuestionResult {
	result := domain.QuestionResult{
		ProblemID:       problem.ID,
		Type:            problem.Type,
		SelectedOptions: evaluateOptions(problem, selectedOptions),
	}
	updateProblemResult(problem, &result, scheme)
	return result
}

func evaluateOptions(problem domain.Problem, selectedOptions []string) []domain.SelectedOption {
	if len(selectedOptions) == 0 {
		return []domain.SelectedOption{}
	}
	options := problem.Options
	if problem.Type == domain.ProblemTypeNumerical {
		isCorrect := false
		if len(options) > 0 {
			isCorrect = options[0].Text == selectedOptions[0]
		}
		return []domain.SelectedOption{
			{
				Text:      selectedOptions[0],
				IsCorrect: isCorrect,
			},
		}
	} else if problem.Type == domain.ProblemTypeSingle {
		selectedOption := domain.SelectedOption{
			ID: selectedOptions[0],
		}
		for _, option := range options {
			if option.ID == selectedOptions[0] {
				selectedOption.Text = option.Text
				selectedOption.IsCorrect = option.IsCorrect
				break
			}
		}
		return []domain.SelectedOption{selectedOption}
	} else {
		res := make([]domain.SelectedOption, 0, len(selectedOptions))
		for _, selectedOptionID := range selectedOptions {
			selectedOption := domain.SelectedOption{ID: selectedOptionID}
			for _, opt := range options {
				if opt.ID == selectedOptionID {
					selectedOption.Text = opt.Text
					selectedOption.IsCorrect = opt.IsCorrect
					break
				}
			}
			res = append(res, selectedOption)
		}
		return res
	}
}

func updateProblemResult(problem domain.Problem, result *domain.QuestionResult, scheme domain.MarkingScheme) {
	correct, incorrect, skipped := marks(result.Type, scheme)
	if problem.Type != domain.ProblemTypeMultiple {
		if len(result.SelectedOptions) == 0 {
			result.Status = domain.QuestionResultStatusSkipped
			result.Marks = skipped
		} else if result.SelectedOptions[0].IsCorrect {
			result.Status = domain.QuestionResultStatusCorrect
			result.Marks = correct
		} else {
			result.Status = domain.QuestionResultStatusIncorrect
			result.Marks = incorrect
		}
	} else {
		if len(result.SelectedOptions) == 0 {
			result.Status = domain.QuestionResultStatusSkipped
			result.Marks = skipped
		} else {
			if hasDuplicateSelectedOptions(result.SelectedOptions) {
				result.Status = domain.QuestionResultStatusIncorrect
				result.Marks = incorrect
				return
			}

			selectedCorrectOptionCount := 0
			for _, option := range result.SelectedOptions {
				if !option.IsCorrect {
					result.Status = domain.QuestionResultStatusIncorrect
					result.Marks = incorrect
					return
				}
				selectedCorrectOptionCount++
			}
			totalCorrectOptionCount := 0
			for _, option := range problem.Options {
				if option.IsCorrect {
					totalCorrectOptionCount++
				}
			}
			if selectedCorrectOptionCount == totalCorrectOptionCount {
				result.Status = domain.QuestionResultStatusCorrect
				result.Marks = correct
			} else if totalCorrectOptionCount > 0 {
				result.Status = domain.QuestionResultStatusPartiallyCorrect
				result.Marks = float64(selectedCorrectOptionCount) / float64(totalCorrectOptionCount) * correct
			} else {
				result.Status = domain.QuestionResultStatusIncorrect
				result.Marks = incorrect
			}
		}
	}
}

func marks(kind domain.ProblemType, s domain.MarkingScheme) (float64, float64, float64) {
	switch kind {
	case domain.ProblemTypeSingle:
		return s.Single.Correct, s.Single.Incorrect, s.Single.Skipped
	case domain.ProblemTypeMultiple:
		return s.Multiple.Correct, s.Multiple.Incorrect, s.Multiple.Skipped
	case domain.ProblemTypeNumerical:
		return s.Numerical.Correct, s.Numerical.Incorrect, s.Numerical.Skipped
	default:
		return 0, 0, 0
	}
}

func hasDuplicateSelectedOptions(options []domain.SelectedOption) bool {
	seen := make(map[string]struct{}, len(options))
	for _, option := range options {
		if _, ok := seen[option.ID]; ok {
			return true
		}
		seen[option.ID] = struct{}{}
	}
	return false
}
