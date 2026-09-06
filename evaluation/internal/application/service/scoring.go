package evaluation

import "github.com/aakashloyar/elevate/evaluation/internal/domain"

func EvaluateAnswer(problem domain.Problem, answer []string, scheme domain.MarkingScheme) domain.QuestionResult {
	result := domain.QuestionResult{
		ProblemID:       problem.ID,
		Type:            problem.Type,
		SelectedOptions: selectedOptions(problem.Type, problem.Options, answer),
	}
	correct, incorrect, skipped := marks(problem.Type, scheme)
	if len(answer) == 0 {
		result.Status, result.Marks = "skipped", skipped
		return result
	}
	if problem.Type == domain.ProblemTypeMultiple {
		return evaluateMultipleAnswer(result, answer, correctOptionIDs(problem.Options), correct, incorrect)
	}
	if sameSet(answer, correctOptionIDs(problem.Options)) {
		result.Status, result.Marks = "correct", correct
	} else {
		result.Status, result.Marks = "incorrect", incorrect
	}
	return result
}

func selectedOptions(problemType domain.ProblemType, options []domain.Option, answer []string) []domain.SelectedOption {
	if len(answer) == 0 {
		return []domain.SelectedOption{}
	}

	if problemType == domain.ProblemTypeNumerical {
		selected := make([]domain.SelectedOption, 0, len(answer))
		for _, value := range answer {
			selected = append(selected, domain.SelectedOption{Text: value})
		}
		return selected
	}

	optionByID := make(map[string]domain.Option, len(options))
	for _, option := range options {
		optionByID[option.ID] = option
	}

	selected := make([]domain.SelectedOption, 0, len(answer))
	for _, optionID := range answer {
		option, ok := optionByID[optionID]
		if !ok {
			continue
		}
		selected = append(selected, domain.SelectedOption{
			ID:        option.ID,
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
		})
	}
	return selected
}

func evaluateMultipleAnswer(result domain.QuestionResult, answer, correctOptionIDs []string, correctMarks, incorrectMarks float64) domain.QuestionResult {
	correctOptions := make(map[string]struct{}, len(correctOptionIDs))
	for _, optionID := range correctOptionIDs {
		correctOptions[optionID] = struct{}{}
	}
	if len(correctOptions) == 0 {
		result.Status, result.Marks = "incorrect", incorrectMarks
		return result
	}

	selected := make(map[string]struct{}, len(answer))
	for _, optionID := range answer {
		if optionID == "" {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
		if _, alreadySelected := selected[optionID]; alreadySelected {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
		selected[optionID] = struct{}{}
		if _, isCorrect := correctOptions[optionID]; !isCorrect {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
	}

	if len(selected) == len(correctOptions) {
		result.Status, result.Marks = "correct", correctMarks
		return result
	}
	result.Status = "partially_correct"
	result.Marks = float64(len(selected)) * correctMarks / float64(len(correctOptions))
	return result
}

func marks(kind domain.ProblemType, s domain.MarkingScheme) (float64, float64, float64) {
	switch kind {
	case domain.ProblemTypeSingle:
		return s.SingleCorrectMarks, s.SingleIncorrectMarks, s.SingleSkippedMarks
	case domain.ProblemTypeMultiple:
		return s.MultipleCorrectMarks, s.MultipleIncorrectMarks, s.MultipleSkippedMarks
	case domain.ProblemTypeNumerical:
		return s.NumericalCorrectMarks, s.NumericalIncorrectMarks, s.NumericalSkippedMarks
	default:
		return 0, 0, 0
	}
}

func correctOptionIDs(options []domain.Option) []string {
	ids := make([]string, 0)
	for _, option := range options {
		if option.IsCorrect {
			ids = append(ids, option.ID)
		}
	}
	return ids
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	values := make(map[string]struct{}, len(a))
	for _, v := range a {
		if v == "" {
			return false
		}
		values[v] = struct{}{}
	}
	if len(values) != len(a) {
		return false
	}
	for _, v := range b {
		if _, ok := values[v]; !ok {
			return false
		}
	}
	return true
}
