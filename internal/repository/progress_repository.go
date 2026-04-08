package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) MarkCompleted(userID, courseID, lessonID int64) error {
	query := `
		INSERT INTO progress (user_id, course_id, lesson_id, completed, completed_at)
		VALUES ($1, $2, $3, TRUE, NOW())
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE SET completed = TRUE, completed_at = NOW()`

	_, err := r.db.Exec(query, userID, courseID, lessonID)
	return err
}

func (r *ProgressRepository) GetCourseProgress(userID, courseID int64) (*models.CourseProgress, error) {
	query := `
		SELECT
			c.id,
			c.title,
			(SELECT COUNT(*) FROM lessons WHERE course_id = c.id) AS total_lessons,
			(SELECT COUNT(*) FROM progress WHERE user_id = $1 AND course_id = c.id AND completed = TRUE) AS completed_count
		FROM courses c
		WHERE c.id = $2`

	cp := &models.CourseProgress{}
	err := r.db.QueryRow(query, userID, courseID).Scan(
		&cp.CourseID, &cp.CourseTitle, &cp.TotalLessons, &cp.CompletedCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course progress: %w", err)
	}

	if cp.TotalLessons > 0 {
		cp.Percentage = float64(cp.CompletedCount) / float64(cp.TotalLessons) * 100
	}

	return cp, nil
}

func (r *ProgressRepository) GetAllProgress(userID int64) ([]models.CourseProgress, error) {
	query := `
		SELECT
			c.id,
			c.title,
			(SELECT COUNT(*) FROM lessons WHERE course_id = c.id) AS total_lessons,
			(SELECT COUNT(*) FROM progress WHERE user_id = $1 AND course_id = c.id AND completed = TRUE) AS completed_count
		FROM courses c
		ORDER BY c.id ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all progress: %w", err)
	}
	defer rows.Close()

	var list []models.CourseProgress
	for rows.Next() {
		var cp models.CourseProgress
		if err := rows.Scan(&cp.CourseID, &cp.CourseTitle, &cp.TotalLessons, &cp.CompletedCount); err != nil {
			return nil, fmt.Errorf("scan progress: %w", err)
		}
		if cp.TotalLessons > 0 {
			cp.Percentage = float64(cp.CompletedCount) / float64(cp.TotalLessons) * 100
		}
		list = append(list, cp)
	}

	return list, rows.Err()
}
