package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/kaplenko/diplom/internal/models"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// pct returns an integer percentage (0-100), rounded to the nearest whole number.
func pct(completed, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(completed) / float64(total) * 100))
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

	cp.Percentage = pct(cp.CompletedCount, cp.TotalLessons)
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
		cp.Percentage = pct(cp.CompletedCount, cp.TotalLessons)
		list = append(list, cp)
	}

	return list, rows.Err()
}

// GetLessonProgress returns task-level completion for a single lesson.
// A task counts as completed when the user has at least one submission with status='passed'.
func (r *ProgressRepository) GetLessonProgress(userID, lessonID int64) (*models.LessonProgress, error) {
	query := `
		SELECT
			l.id,
			l.title,
			(SELECT COUNT(*) FROM tasks WHERE lesson_id = l.id) AS total_tasks,
			(SELECT COUNT(DISTINCT t.id)
			   FROM tasks t
			   JOIN submissions s ON s.task_id = t.id AND s.user_id = $1 AND s.status = 'passed'
			  WHERE t.lesson_id = l.id
			) AS completed_tasks
		FROM lessons l
		WHERE l.id = $2`

	lp := &models.LessonProgress{}
	err := r.db.QueryRow(query, userID, lessonID).Scan(
		&lp.LessonID, &lp.LessonTitle, &lp.TotalTasks, &lp.CompletedTasks,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson progress: %w", err)
	}

	lp.Percentage = pct(lp.CompletedTasks, lp.TotalTasks)
	return lp, nil
}

// GetCourseLessonProgress returns per-lesson task completion for every
// lesson in a course — useful for rendering a detailed progress breakdown.
func (r *ProgressRepository) GetCourseLessonProgress(userID, courseID int64) ([]models.LessonProgress, error) {
	query := `
		SELECT
			l.id,
			l.title,
			(SELECT COUNT(*) FROM tasks WHERE lesson_id = l.id) AS total_tasks,
			(SELECT COUNT(DISTINCT t.id)
			   FROM tasks t
			   JOIN submissions s ON s.task_id = t.id AND s.user_id = $1 AND s.status = 'passed'
			  WHERE t.lesson_id = l.id
			) AS completed_tasks
		FROM lessons l
		WHERE l.course_id = $2
		ORDER BY l.order_index ASC`

	rows, err := r.db.Query(query, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course lesson progress: %w", err)
	}
	defer rows.Close()

	var list []models.LessonProgress
	for rows.Next() {
		var lp models.LessonProgress
		if err := rows.Scan(&lp.LessonID, &lp.LessonTitle, &lp.TotalTasks, &lp.CompletedTasks); err != nil {
			return nil, fmt.Errorf("scan lesson progress: %w", err)
		}
		lp.Percentage = pct(lp.CompletedTasks, lp.TotalTasks)
		list = append(list, lp)
	}

	return list, rows.Err()
}
