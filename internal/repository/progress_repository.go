package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/doug-martin/goqu/v9"

	"github.com/kaplenko/diplom/internal/models"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func pct(completed, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(completed) / float64(total) * 100))
}

func (r *ProgressRepository) MarkCompleted(userID, courseID, lessonID int64) error {
	query, args, err := pg.Insert("progress").
		Cols("user_id", "course_id", "lesson_id", "completed", "completed_at").
		Vals(goqu.Vals{userID, courseID, lessonID, true, goqu.L("NOW()")}).
		OnConflict(goqu.DoUpdate("user_id, lesson_id", goqu.Record{
			"completed":    true,
			"completed_at": goqu.L("NOW()"),
		})).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build upsert query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *ProgressRepository) GetCourseProgress(userID, courseID int64) (*models.CourseProgress, error) {
	totalLessons := pg.From("lessons").
		Select(goqu.COUNT(goqu.Star())).
		Where(goqu.C("course_id").Eq(goqu.I("c.id")))

	completedCount := pg.From("progress").
		Select(goqu.COUNT(goqu.Star())).
		Where(
			goqu.C("user_id").Eq(userID),
			goqu.C("course_id").Eq(goqu.I("c.id")),
			goqu.C("completed").Eq(true),
		)

	query, args, err := pg.From(goqu.T("courses").As("c")).
		Select(
			goqu.I("c.id"),
			goqu.I("c.title"),
			totalLessons.As("total_lessons"),
			completedCount.As("completed_count"),
		).
		Where(goqu.I("c.id").Eq(courseID)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build course progress query: %w", err)
	}

	cp := &models.CourseProgress{}
	err = r.db.QueryRow(query, args...).Scan(
		&cp.CourseID, &cp.CourseTitle, &cp.TotalLessons, &cp.CompletedCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course progress: %w", err)
	}

	cp.Percentage = pct(cp.CompletedCount, cp.TotalLessons)
	return cp, nil
}

func (r *ProgressRepository) GetAllProgress(userID int64) ([]models.CourseProgress, error) {
	totalLessons := pg.From("lessons").
		Select(goqu.COUNT(goqu.Star())).
		Where(goqu.C("course_id").Eq(goqu.I("c.id")))

	completedCount := pg.From("progress").
		Select(goqu.COUNT(goqu.Star())).
		Where(
			goqu.C("user_id").Eq(userID),
			goqu.C("course_id").Eq(goqu.I("c.id")),
			goqu.C("completed").Eq(true),
		)

	query, args, err := pg.From(goqu.T("courses").As("c")).
		Select(
			goqu.I("c.id"),
			goqu.I("c.title"),
			totalLessons.As("total_lessons"),
			completedCount.As("completed_count"),
		).
		Order(goqu.I("c.id").Asc()).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build all progress query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
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
	totalTasks := pg.From("tasks").
		Select(goqu.COUNT(goqu.Star())).
		Where(goqu.C("lesson_id").Eq(goqu.I("l.id")))

	completedTasks := pg.From(goqu.T("tasks").As("t")).
		Join(
			goqu.T("submissions").As("s"),
			goqu.On(
				goqu.I("s.task_id").Eq(goqu.I("t.id")),
				goqu.I("s.user_id").Eq(userID),
				goqu.I("s.status").Eq("passed"),
			),
		).
		Select(goqu.L("COUNT(DISTINCT ?)", goqu.I("t.id"))).
		Where(goqu.I("t.lesson_id").Eq(goqu.I("l.id")))

	query, args, err := pg.From(goqu.T("lessons").As("l")).
		Select(
			goqu.I("l.id"),
			goqu.I("l.title"),
			totalTasks.As("total_tasks"),
			completedTasks.As("completed_tasks"),
		).
		Where(goqu.I("l.id").Eq(lessonID)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build lesson progress query: %w", err)
	}

	lp := &models.LessonProgress{}
	err = r.db.QueryRow(query, args...).Scan(
		&lp.LessonID, &lp.LessonTitle, &lp.TotalTasks, &lp.CompletedTasks,
	)
	if errors.Is(err, sql.ErrNoRows) {
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
	totalTasks := pg.From("tasks").
		Select(goqu.COUNT(goqu.Star())).
		Where(goqu.C("lesson_id").Eq(goqu.I("l.id")))

	completedTasks := pg.From(goqu.T("tasks").As("t")).
		Join(
			goqu.T("submissions").As("s"),
			goqu.On(
				goqu.I("s.task_id").Eq(goqu.I("t.id")),
				goqu.I("s.user_id").Eq(userID),
				goqu.I("s.status").Eq("passed"),
			),
		).
		Select(goqu.L("COUNT(DISTINCT ?)", goqu.I("t.id"))).
		Where(goqu.I("t.lesson_id").Eq(goqu.I("l.id")))

	query, args, err := pg.From(goqu.T("lessons").As("l")).
		Select(
			goqu.I("l.id"),
			goqu.I("l.title"),
			totalTasks.As("total_tasks"),
			completedTasks.As("completed_tasks"),
		).
		Where(goqu.I("l.course_id").Eq(courseID)).
		Order(goqu.I("l.order_index").Asc()).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build course lesson progress query: %w", err)
	}

	rows, err := r.db.Query(query, args...)
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
