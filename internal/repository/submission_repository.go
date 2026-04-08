package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/kaplenko/diplom/internal/models"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(sub *models.Submission) error {
	query, args, err := pg.Insert("submissions").
		Cols("task_id", "user_id", "code", "status").
		Vals(goqu.Vals{sub.TaskID, sub.UserID, sub.Code, sub.Status}).
		Returning("id", "submitted_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&sub.ID, &sub.SubmittedAt)
}

func (r *SubmissionRepository) GetByID(id int64) (*models.Submission, error) {
	s := &models.Submission{}

	query, args, err := pg.From("submissions").
		Select("id", "task_id", "user_id", "code", "status", "result", "score", "submitted_at").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&s.ID, &s.TaskID, &s.UserID, &s.Code, &s.Status,
		&s.Result, &s.Score, &s.SubmittedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get submission by id: %w", err)
	}
	return s, nil
}

func (r *SubmissionRepository) UpdateStatus(id int64, status models.SubmissionStatus, result string, score int) error {
	query, args, err := pg.Update("submissions").
		Set(goqu.Record{
			"status": status,
			"result": result,
			"score":  score,
		}).
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *SubmissionRepository) ListByTask(taskID, userID int64, params models.PaginationParams) ([]models.Submission, int64, error) {
	var total int64

	base := pg.From("submissions").Where(
		goqu.C("task_id").Eq(taskID),
		goqu.C("user_id").Eq(userID),
	)

	countSQL, countArgs, err := base.Select(goqu.COUNT(goqu.Star())).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count submissions: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "task_id", "user_id", "code", "status", "result", "score", "submitted_at").
		Order(goqu.C("submitted_at").Desc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
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

	base := pg.From("submissions")

	countSQL, countArgs, err := base.Select(goqu.COUNT(goqu.Star())).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count all submissions: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "task_id", "user_id", "code", "status", "result", "score", "submitted_at").
		Order(goqu.C("submitted_at").Desc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
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
