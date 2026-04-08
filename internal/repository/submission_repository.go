package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(sub *models.Submission) error {
	query := `
		INSERT INTO submissions (task_id, user_id, code, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, submitted_at`

	return r.db.QueryRow(query, sub.TaskID, sub.UserID, sub.Code, sub.Status).
		Scan(&sub.ID, &sub.SubmittedAt)
}

func (r *SubmissionRepository) GetByID(id int64) (*models.Submission, error) {
	s := &models.Submission{}
	query := `
		SELECT id, task_id, user_id, code, status, result, score, submitted_at
		FROM submissions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.TaskID, &s.UserID, &s.Code, &s.Status,
		&s.Result, &s.Score, &s.SubmittedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get submission by id: %w", err)
	}
	return s, nil
}

func (r *SubmissionRepository) UpdateStatus(id int64, status models.SubmissionStatus, result string, score int) error {
	query := `UPDATE submissions SET status = $1, result = $2, score = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, result, score, id)
	return err
}

func (r *SubmissionRepository) ListByTask(taskID, userID int64, params models.PaginationParams) ([]models.Submission, int64, error) {
	var total int64

	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM submissions WHERE task_id = $1 AND user_id = $2`,
		taskID, userID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count submissions: %w", err)
	}

	query := `
		SELECT id, task_id, user_id, code, status, result, score, submitted_at
		FROM submissions WHERE task_id = $1 AND user_id = $2
		ORDER BY submitted_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, taskID, userID, params.PageSize, params.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("list submissions: %w", err)
	}
	defer rows.Close()

	var subs []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.TaskID, &s.UserID, &s.Code, &s.Status, &s.Result, &s.Score, &s.SubmittedAt); err != nil {
			return nil, 0, fmt.Errorf("scan submission: %w", err)
		}
		subs = append(subs, s)
	}

	return subs, total, rows.Err()
}

func (r *SubmissionRepository) ListAll(params models.PaginationParams) ([]models.Submission, int64, error) {
	var total int64

	if err := r.db.QueryRow(`SELECT COUNT(*) FROM submissions`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count all submissions: %w", err)
	}

	query := `
		SELECT id, task_id, user_id, code, status, result, score, submitted_at
		FROM submissions
		ORDER BY submitted_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, params.PageSize, params.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("list all submissions: %w", err)
	}
	defer rows.Close()

	var subs []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.TaskID, &s.UserID, &s.Code, &s.Status, &s.Result, &s.Score, &s.SubmittedAt); err != nil {
			return nil, 0, fmt.Errorf("scan submission: %w", err)
		}
		subs = append(subs, s)
	}

	return subs, total, rows.Err()
}
