package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/kaplenko/diplom/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query, args, err := pg.Insert("users").
		Cols("email", "password_hash", "name", "role").
		Vals(goqu.Vals{user.Email, user.PasswordHash, user.Name, user.Role}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}

	query, args, err := pg.From("users").
		Select("id", "email", "password_hash", "name", "role", "created_at", "updated_at").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}

	query, args, err := pg.From("users").
		Select("id", "email", "password_hash", "name", "role", "created_at", "updated_at").
		Where(goqu.C("email").Eq(email)).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	err = r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query, args, err := pg.Update("users").
		Set(goqu.Record{
			"name":       user.Name,
			"updated_at": goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(user.ID)).
		Returning("updated_at").
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	return r.db.QueryRow(query, args...).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(id int64) error {
	query, args, err := pg.Delete("users").
		Where(goqu.C("id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.db.Exec(query, args...)
	return err
}

func (r *UserRepository) List(params models.PaginationParams) ([]models.User, int64, error) {
	var total int64

	base := pg.From("users")
	if params.Search != "" {
		pattern := "%" + params.Search + "%"
		base = base.Where(goqu.Or(
			goqu.C("name").ILike(pattern),
			goqu.C("email").ILike(pattern),
		))
	}

	countSQL, countArgs, err := base.Select(goqu.COUNT(goqu.Star())).Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}
	if err := r.db.QueryRow(countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	listSQL, listArgs, err := base.
		Select("id", "email", "password_hash", "name", "role", "created_at", "updated_at").
		Order(goqu.C("id").Asc()).
		Limit(uint(params.PageSize)).
		Offset(uint(params.Offset())).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.db.Query(listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, total, rows.Err()
}
