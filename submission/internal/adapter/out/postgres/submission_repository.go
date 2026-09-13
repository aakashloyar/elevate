package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) out.SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Save(submission domain.Submission, drafts []domain.SubmissionAnswerDraft) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
		INSERT INTO submissions (
			id,
			assessment_id,
			user_id,
			status,
			started_at,
			duration_seconds,
			expires_at,
			submitted_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	if _, err := tx.Exec(query, submission.ID, submission.AssessmentID, submission.UserID, submission.Status, submission.StartedAt, submission.DurationSeconds, submission.ExpiresAt, submission.SubmittedAt, submission.CreatedAt, submission.UpdatedAt); err != nil {
		return err
	}
	for _, draft := range drafts {
		optionIDs := make([]string, 0, len(draft.Options))
		optionTexts := make([]string, 0, len(draft.Options))
		for _, option := range draft.Options {
			optionIDs = append(optionIDs, option.ID)
			optionTexts = append(optionTexts, option.Text)
		}
		if _, err := tx.Exec(`INSERT INTO submission_answer_drafts (submission_id, problem_id, problem_type, option_ids, option_texts, answer, answer_updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, submission.ID, draft.ProblemID, draft.ProblemType, pqStringArray(optionIDs), pqStringArray(optionTexts), pqStringArray(draft.Answer), draft.AnswerUpdatedAt); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SubmissionRepository) SaveAnswer(answer domain.SubmissionAnswerDraft) (bool, error) {
	updatedAt := time.Now().UTC()
	if answer.AnswerUpdatedAt != nil {
		updatedAt = *answer.AnswerUpdatedAt
	}
	answerUpdatedAt := updatedAt
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var status domain.SubmissionStatus
	if err := tx.QueryRow(`SELECT status FROM submissions WHERE id = $1 FOR UPDATE`, answer.SubmissionID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if status != domain.SubmissionStatusInProgress {
		return false, nil
	}

	result, err := tx.Exec(`
		UPDATE submission_answer_drafts
		SET answer = $3,
			answer_updated_at = $4
		WHERE submission_id = $1
			AND problem_id = $2
	`, answer.SubmissionID, answer.ProblemID, pqStringArray(answer.Answer), answerUpdatedAt)
	if err != nil {
		return false, err
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if updated == 0 {
		return false, nil
	}
	if _, err := tx.Exec(`UPDATE submissions SET updated_at = $2 WHERE id = $1`, answer.SubmissionID, updatedAt); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *SubmissionRepository) SaveAnswers(answers []domain.SubmissionAnswerDraft) (bool, error) {
	if len(answers) == 0 {
		return true, nil
	}

	submissionID := answers[0].SubmissionID
	updatedAt := time.Now().UTC()
	if answers[0].AnswerUpdatedAt != nil {
		updatedAt = *answers[0].AnswerUpdatedAt
	}

	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var status domain.SubmissionStatus
	if err := tx.QueryRow(`SELECT status FROM submissions WHERE id = $1 FOR UPDATE`, submissionID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if status != domain.SubmissionStatusInProgress {
		return false, nil
	}

	for _, answer := range answers {
		answerUpdatedAt := updatedAt
		if answer.AnswerUpdatedAt != nil {
			answerUpdatedAt = *answer.AnswerUpdatedAt
		}
		result, err := tx.Exec(`
			UPDATE submission_answer_drafts
			SET answer = $3,
				answer_updated_at = $4
			WHERE submission_id = $1
				AND problem_id = $2
		`, answer.SubmissionID, answer.ProblemID, pqStringArray(answer.Answer), answerUpdatedAt)
		if err != nil {
			return false, err
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		if updated == 0 {
			return false, nil
		}
		if answerUpdatedAt.After(updatedAt) {
			updatedAt = answerUpdatedAt
		}
	}
	if _, err := tx.Exec(`UPDATE submissions SET updated_at = $2 WHERE id = $1`, submissionID, updatedAt); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *SubmissionRepository) FindAnswerSnapshots(submissionID string) (domain.SubmissionStatus, []domain.SubmissionAnswerDraft, error) {
	rows, err := r.db.Query(`
		SELECT
			submission.status,
			draft.problem_id,
			draft.problem_type,
			draft.option_ids,
			draft.option_texts,
			draft.answer,
			draft.answer_updated_at
		FROM submissions AS submission
		INNER JOIN submission_answer_drafts AS draft
			ON draft.submission_id = submission.id
		WHERE submission.id = $1
		ORDER BY draft.problem_id
	`, submissionID)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var status domain.SubmissionStatus
	drafts := make([]domain.SubmissionAnswerDraft, 0)
	found := false
	for rows.Next() {
		found = true
		var draft domain.SubmissionAnswerDraft
		var optionIDs []string
		var optionTexts []string
		var answer []string
		var answerUpdatedAt sql.NullTime
		if err := rows.Scan(&status, &draft.ProblemID, &draft.ProblemType, pq.Array(&optionIDs), pq.Array(&optionTexts), pq.Array(&answer), &answerUpdatedAt); err != nil {
			return "", nil, err
		}
		draft.SubmissionID = submissionID
		draft.Answer = answer
		for index, optionID := range optionIDs {
			option := domain.Option{ID: optionID}
			if index < len(optionTexts) {
				option.Text = optionTexts[index]
			}
			draft.Options = append(draft.Options, option)
		}
		if answerUpdatedAt.Valid {
			draft.AnswerUpdatedAt = &answerUpdatedAt.Time
		}
		drafts = append(drafts, draft)
	}
	if err := rows.Err(); err != nil {
		return "", nil, err
	}
	if !found {
		return "", nil, sql.ErrNoRows
	}
	return status, drafts, nil
}

func (r *SubmissionRepository) FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswerDraft, error) {
	submissionQuery := `SELECT id, assessment_id, user_id, status, started_at, duration_seconds, expires_at, submitted_at, created_at, updated_at FROM submissions WHERE id = $1`
	row := r.db.QueryRow(submissionQuery, submissionID)

	var submission domain.Submission
	var startedAt sql.NullTime
	var submittedAt sql.NullTime
	var expiresAt sql.NullTime
	if err := row.Scan(&submission.ID, &submission.AssessmentID, &submission.UserID, &submission.Status, &startedAt, &submission.DurationSeconds, &expiresAt, &submittedAt, &submission.CreatedAt, &submission.UpdatedAt); err != nil {
		return domain.Submission{}, nil, err
	}
	if startedAt.Valid {
		submission.StartedAt = &startedAt.Time
	}
	if expiresAt.Valid {
		submission.ExpiresAt = &expiresAt.Time
	}
	if submittedAt.Valid {
		submission.SubmittedAt = &submittedAt.Time
	}
	draftRows, err := r.db.Query(`SELECT problem_id, problem_type, option_ids, option_texts, answer, answer_updated_at FROM submission_answer_drafts WHERE submission_id = $1 ORDER BY problem_id`, submissionID)
	if err != nil {
		return domain.Submission{}, nil, err
	}
	defer draftRows.Close()
	drafts := make([]domain.SubmissionAnswerDraft, 0)
	for draftRows.Next() {
		var draft domain.SubmissionAnswerDraft
		var optionIDs []string
		var optionTexts []string
		var answer []string
		var answerUpdatedAt sql.NullTime
		if err := draftRows.Scan(&draft.ProblemID, &draft.ProblemType, pq.Array(&optionIDs), pq.Array(&optionTexts), pq.Array(&answer), &answerUpdatedAt); err != nil {
			return domain.Submission{}, nil, err
		}
		draft.SubmissionID = submissionID
		draft.Answer = answer
		for index, optionID := range optionIDs {
			option := domain.Option{ID: optionID}
			if index < len(optionTexts) {
				option.Text = optionTexts[index]
			}
			draft.Options = append(draft.Options, option)
		}
		if answerUpdatedAt.Valid {
			draft.AnswerUpdatedAt = &answerUpdatedAt.Time
		}
		drafts = append(drafts, draft)
	}
	if err := draftRows.Err(); err != nil {
		return domain.Submission{}, nil, err
	}
	return submission, drafts, nil
}

func (r *SubmissionRepository) FindStatus(submissionID string) (domain.SubmissionStatus, *time.Time, error) {
	var status domain.SubmissionStatus
	var expiresAt sql.NullTime
	if err := r.db.QueryRow(`SELECT status, expires_at FROM submissions WHERE id = $1`, submissionID).Scan(&status, &expiresAt); err != nil {
		return "", nil, err
	}
	if expiresAt.Valid {
		return status, &expiresAt.Time, nil
	}
	return status, nil, nil
}

func (r *SubmissionRepository) UpdateStatus(submissionID string, status domain.SubmissionStatus) error {
	_, err := r.db.Exec(`UPDATE submissions SET status = $2, updated_at = NOW() WHERE id = $1`, submissionID, status)
	return err
}

func (r *SubmissionRepository) UpdateStartTime(submissionID string, startedAt, expiresAt time.Time, status domain.SubmissionStatus) error {
	_, err := r.db.Exec(`UPDATE submissions SET started_at = $2, expires_at = $3, status = $4, updated_at = NOW() WHERE id = $1 AND status = 'CREATED'`, submissionID, startedAt, expiresAt, status)
	return err
}

func (r *SubmissionRepository) Submit(submissionID string, submittedAt time.Time) (bool, error) {
	result, err := r.db.Exec(`UPDATE submissions
		SET status = 'SUBMITTED', submitted_at = $2, updated_at = $2
		WHERE id = $1 AND status = 'IN_PROGRESS' AND expires_at >= $2`, submissionID, submittedAt)
	if err != nil {
		return false, err
	}
	updated, err := result.RowsAffected()
	return updated == 1, err
}

func (r *SubmissionRepository) ExpireSubmissions(expiredAt time.Time, limit int) ([]string, error) {
	rows, err := r.db.Query(`WITH expired AS (
		SELECT id
		FROM submissions
		WHERE status = 'IN_PROGRESS' AND expires_at <= $1
		ORDER BY expires_at
		FOR UPDATE SKIP LOCKED
		LIMIT $2
	)
	UPDATE submissions AS submission
	SET status = 'SUBMITTED', submitted_at = $1, updated_at = $1
	FROM expired
	WHERE submission.id = expired.id AND submission.status = 'IN_PROGRESS'
	RETURNING submission.id`, expiredAt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissionIDs := make([]string, 0)
	for rows.Next() {
		var submissionID string
		if err := rows.Scan(&submissionID); err != nil {
			return nil, err
		}
		submissionIDs = append(submissionIDs, submissionID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return submissionIDs, nil
}

func pqStringArray(values []string) []string {
	return values
}
