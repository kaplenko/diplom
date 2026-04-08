package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *models.Course) error {
	query := `
		INSERT INTO courses (title, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, course.Title, course.Description, course.CreatedBy).
		Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
}

func (r *CourseRepository) GetByID(id int64) (*models.Course, error) {
	c := &models.Course{}
	query := `SELECT id, title, description, created_by, created_at, updated_at FROM courses WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Title, &c.Description, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course by id: %w", err)
	}
	return c, nil
}

func (r *CourseRepository) Update(course *models.Course) error {
	query := `
		UPDATE courses SET title = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	return r.db.QueryRow(query, course.Title, course.Description, course.ID).Scan(&course.UpdatedAt)
}

func (r *CourseRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM courses WHERE id = $1`, id)
	return err
}

func (r *CourseRepository) List(params models.PaginationParams) ([]models.Course, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM courses`
	listQuery := `SELECT id, title, description, created_by, created_at, updated_at FROM courses`

	args := []interface{}{}
	where := ""
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(` WHERE title ILIKE $%d OR description ILIKE $%d`, argIdx, argIdx+1)
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	if err := r.db.QueryRow(countQuery+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count courses: %w", err)
	}

	listQuery += where + fmt.Sprintf(` ORDER BY id ASC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, params.PageSize, params.Offset())

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan course: %w", err)
		}
		courses = append(courses, c)
	}

	return courses, total, rows.Err()
}
