package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/kaplenko/diplom/internal/models"
)

var pg = goqu.Dialect("postgres")

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *models.Course) error {
	query, args, err := pg.Insert("courses").
		Cols("title", "description", "created_by").
		Vals(goqu.Vals{course.Title, course.Description, course.CreatedBy}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
}

func (r *CourseRepository) GetByID(id int64) (*models.Course, error) {
	c := &models.Course{}

	query, args, err := pg.From("courses").
		Select("id", "title", "description", "created_by", "created_at", "updated_at").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&c.ID, &c.Title, &c.Description, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course by id: %w", err)
	}
	return c, nil
}

func (r *CourseRepository) Update(course *models.Course) error {
	query, args, err := pg.Update("courses").
		Set(goqu.Record{
			"title":       course.Title,
			"description": course.Description,
			"updated_at":  goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(course.ID)).
		Returning("updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&course.UpdatedAt)
}

func (r *CourseRepository) Delete(id int64) error {
	query, args, err := pg.Delete("courses").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *CourseRepository) List(params models.PaginationParams) ([]models.Course, int64, error) {
	var total int64

	base := pg.From("courses")
	if params.Search != "" {
		pattern := "%" + params.Search + "%"
		base = base.Where(goqu.Or(
			goqu.C("title").ILike(pattern),
			goqu.C("description").ILike(pattern),
		))
	}

	countSQL, countArgs, err := base.Select(goqu.COUNT(goqu.Star())).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count courses: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "title", "description", "created_by", "created_at", "updated_at").
		Order(goqu.C("id").Asc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
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
