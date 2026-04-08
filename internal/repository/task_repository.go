package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (lesson_id, title, description, initial_code, test_cases, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		task.LessonID, task.Title, task.Description,
		task.InitialCode, task.TestCases, task.Difficulty,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	t := &models.Task{}
	query := `
		SELECT id, lesson_id, title, description, initial_code, test_cases, difficulty, created_at, updated_at
		FROM tasks WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.LessonID, &t.Title, &t.Description,
		&t.InitialCode, &t.TestCases, &t.Difficulty,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return t, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, initial_code = $3, test_cases = $4, difficulty = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(query,
		task.Title, task.Description, task.InitialCode,
		task.TestCases, task.Difficulty, task.ID,
	).Scan(&task.UpdatedAt)
}

func (r *TaskRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *TaskRepository) ListByLesson(lessonID int64, params models.PaginationParams) ([]models.Task, int64, error) {
	var total int64

	if err := r.db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE lesson_id = $1`, lessonID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	query := `
		SELECT id, lesson_id, title, description, initial_code, test_cases, difficulty, created_at, updated_at
		FROM tasks WHERE lesson_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, lessonID, params.PageSize, params.Offset())
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
