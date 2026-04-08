package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/kaplenko/diplom/internal/models"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(lesson *models.Lesson) error {
	query, args, err := pg.Insert("lessons").
		Cols("course_id", "title", "content", "order_index").
		Vals(goqu.Vals{lesson.CourseID, lesson.Title, lesson.Content, lesson.OrderIndex}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&lesson.ID, &lesson.CreatedAt, &lesson.UpdatedAt)
}

func (r *LessonRepository) GetByID(id int64) (*models.Lesson, error) {
	l := &models.Lesson{}

	query, args, err := pg.From("lessons").
		Select("id", "course_id", "title", "content", "order_index", "created_at", "updated_at").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Content, &l.OrderIndex, &l.CreatedAt, &l.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson by id: %w", err)
	}
	return l, nil
}

func (r *LessonRepository) Update(lesson *models.Lesson) error {
	query, args, err := pg.Update("lessons").
		Set(goqu.Record{
			"title":       lesson.Title,
			"content":     lesson.Content,
			"order_index": lesson.OrderIndex,
			"updated_at":  goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(lesson.ID)).
		Returning("updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&lesson.UpdatedAt)
}

func (r *LessonRepository) Delete(id int64) error {
	query, args, err := pg.
		From("lessons").
		Delete().
		Where(goqu.C("id").Eq(id)).
		Returning("updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *LessonRepository) ListByCourse(courseID int64, params models.PaginationParams) ([]models.Lesson, int64, error) {
	var total int64

	base := pg.From("lessons").Where(goqu.C("course_id").Eq(courseID))

	countSQL, countArgs, err := base.
		Select(goqu.COUNT(goqu.Star())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count lessons: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "course_id", "title", "content", "order_index", "created_at", "updated_at").
		Order(goqu.C("order_index").Asc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
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
