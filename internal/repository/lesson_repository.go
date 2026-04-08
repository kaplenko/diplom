package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (course_id, title, content, order_index)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, lesson.CourseID, lesson.Title, lesson.Content, lesson.OrderIndex).
		Scan(&lesson.ID, &lesson.CreatedAt, &lesson.UpdatedAt)
}

func (r *LessonRepository) GetByID(id int64) (*models.Lesson, error) {
	l := &models.Lesson{}
	query := `SELECT id, course_id, title, content, order_index, created_at, updated_at FROM lessons WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson by id: %w", err)
	}
	return l, nil
}

func (r *LessonRepository) Update(lesson *models.Lesson) error {
	query := `
		UPDATE lessons SET title = $1, content = $2, order_index = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	return r.db.QueryRow(query, lesson.Title, lesson.Content, lesson.OrderIndex, lesson.ID).
		Scan(&lesson.UpdatedAt)
}

func (r *LessonRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM lessons WHERE id = $1`, id)
	return err
}

func (r *LessonRepository) ListByCourse(courseID int64, params models.PaginationParams) ([]models.Lesson, int64, error) {
	var total int64

	if err := r.db.QueryRow(`SELECT COUNT(*) FROM lessons WHERE course_id = $1`, courseID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count lessons: %w", err)
	}

	query := `
		SELECT id, course_id, title, content, order_index, created_at, updated_at
		FROM lessons WHERE course_id = $1
		ORDER BY order_index ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, courseID, params.PageSize, params.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("list lessons: %w", err)
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		if err := rows.Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan lesson: %w", err)
		}
		lessons = append(lessons, l)
	}

	return lessons, total, rows.Err()
}
