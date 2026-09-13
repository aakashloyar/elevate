package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

const sourceTypeAI = "ai"

type AIProblemGenerator struct {
	ai        out.AITextGenerator
	publisher out.GeneratedProblemPublisher
	topic     string
}

func NewAIProblemGenerator(ai out.AITextGenerator, publisher out.GeneratedProblemPublisher, topic string) out.GenerationProcessor {
	return &AIProblemGenerator{ai: ai, publisher: publisher, topic: topic}
}

func (p *AIProblemGenerator) ProcessGeneration(ctx context.Context, job domain.GenerationJob) error {
	if job.AssessmentID == nil || strings.TrimSpace(*job.AssessmentID) == "" {
		return errors.New("assessment id is required to publish generated problems")
	}

	response, err := p.ai.GenerateText(ctx, out.GenerateTextRequest{
		SystemPrompt:     systemPrompt(),
		UserPrompt:       userPrompt(job),
		ResponseMimeType: "application/json",
		Temperature:      floatPtr(0.2),
	})
	if err != nil {
		return fmt.Errorf("generate problems with ai: %w", err)
	}

	batch, err := parseGeneratedProblemBatch(response.Text)
	if err != nil {
		return err
	}

	problems := normalizeGeneratedProblems(job, batch.Problems)
	if err := validateGeneratedProblems(problems); err != nil {
		return err
	}

	return p.publisher.PublishGeneratedProblems(ctx, out.GeneratedProblemBatchMessage{
		Topic: p.topic,
		Event: out.GeneratedProblemBatchEvent{
			AssessmentID: strings.TrimSpace(*job.AssessmentID),
			Problems:     problems,
		},
	})
}

type generatedProblemBatch struct {
	Problems []out.GeneratedProblem `json:"problems"`
}

func systemPrompt() string {
	return strings.Join([]string{
		"You generate assessment problems as strict JSON only.",
		"Return exactly one JSON object with a problems array.",
		"Each problem must include title, statement, type, difficulty, options, and tags.",
		"type must be one of: single, multiple, numerical.",
		"difficulty must be one of: easy, medium, hard.",
		"For single problems, exactly one option must be correct.",
		"For multiple problems, at least one option must be correct.",
		"For numerical problems, include exactly one correct option whose text is the answer.",
		"Do not include markdown fences or explanation.",
	}, " ")
}

func userPrompt(job domain.GenerationJob) string {
	return fmt.Sprintf(
		"Create %d single-correct, %d multiple-correct, and %d numerical problems. Level: %s. Description: %s. Topic IDs: %s.",
		job.SingleCorrectCount,
		job.MultiCorrectCount,
		job.NumericalCount,
		job.Level,
		job.Description,
		strings.Join(job.TopicIDs, ", "),
	)
}

func parseGeneratedProblemBatch(text string) (generatedProblemBatch, error) {
	var batch generatedProblemBatch
	if err := json.Unmarshal([]byte(cleanJSON(text)), &batch); err != nil {
		return generatedProblemBatch{}, fmt.Errorf("parse generated problems: %w", err)
	}
	return batch, nil
}

func cleanJSON(text string) string {
	trimmed := strings.TrimSpace(text)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	return strings.TrimSpace(trimmed)
}

func validateGeneratedProblems(problems []out.GeneratedProblem) error {
	if len(problems) == 0 {
		return errors.New("ai response did not include any problems")
	}

	for index, problem := range problems {
		if strings.TrimSpace(problem.Statement) == "" {
			return fmt.Errorf("problem %d statement is required", index)
		}
		if !problem.Type.IsValid() {
			return fmt.Errorf("problem %d has invalid type %q", index, problem.Type)
		}
		if !problem.Difficulty.IsValid() {
			return fmt.Errorf("problem %d has invalid difficulty %q", index, problem.Difficulty)
		}
		if len(problem.Options) == 0 {
			return fmt.Errorf("problem %d options are required", index)
		}
		correctCount := 0
		for _, option := range problem.Options {
			if strings.TrimSpace(option.Text) == "" {
				return fmt.Errorf("problem %d option text is required", index)
			}
			if option.IsCorrect {
				correctCount++
			}
		}
		switch problem.Type {
		case "single", "numerical":
			if correctCount != 1 {
				return fmt.Errorf("problem %d must have exactly one correct option", index)
			}
		case "multiple":
			if correctCount < 1 {
				return fmt.Errorf("problem %d must have at least one correct option", index)
			}
		}
	}
	return nil
}

func normalizeGeneratedProblems(job domain.GenerationJob, problems []out.GeneratedProblem) []out.GeneratedProblem {
	normalized := make([]out.GeneratedProblem, 0, len(problems))
	for _, problem := range problems {
		problem.CreatedBy = strings.TrimSpace(job.UserID)
		problem.SourceType = sourceTypeAI
		if strings.TrimSpace(string(problem.Difficulty)) == "" {
			problem.Difficulty = job.Level
		}
		normalized = append(normalized, problem)
	}
	return normalized
}

func floatPtr(value float64) *float64 {
	return &value
}
