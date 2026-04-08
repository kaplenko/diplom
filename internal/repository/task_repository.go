package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/kaplenko/diplom/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query, args, err := pg.Insert("tasks").
		Cols("lesson_id", "title", "description", "initial_code", "test_cases", "difficulty").
		Vals(goqu.Vals{task.LessonID, task.Title, task.Description, task.InitialCode, string(task.TestCases), task.Difficulty}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	t := &models.Task{}

	query, args, err := pg.From("tasks").
		Select("id", "lesson_id", "title", "description", "initial_code", "test_cases", "difficulty", "created_at", "updated_at").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&t.ID, &t.LessonID, &t.Title, &t.Description,
		&t.InitialCode, &t.TestCases, &t.Difficulty,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return t, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	query, args, err := pg.Update("tasks").
		Set(goqu.Record{
			"title":        task.Title,
			"description":  task.Description,
			"initial_code": task.InitialCode,
			"test_cases":   string(task.TestCases),
			"difficulty":   task.Difficulty,
			"updated_at":   goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(task.ID)).
		Returning("updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&task.UpdatedAt)
}

func (r *TaskRepository) Delete(id int64) error {
	query, args, err := pg.Delete("tasks").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *TaskRepository) ListByLesson(lessonID int64, params models.PaginationParams) ([]models.Task, int64, error) {
	var total int64

	base := pg.From("tasks").Where(goqu.C("lesson_id").Eq(lessonID))

	countSQL, countArgs, err := base.Select(goqu.COUNT(goqu.Star())).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "lesson_id", "title", "description", "initial_code", "test_cases", "difficulty", "created_at", "updated_at").
		Order(goqu.C("id").Asc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID, &t.LessonID, &t.Title, &t.Description,
			&t.InitialCode, &t.TestCases, &t.Difficulty,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, total, rows.Err()
}
