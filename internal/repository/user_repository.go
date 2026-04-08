package repository

import (
	"database/sql"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET name = $1, updated_at = NOW() WHERE id = $2 RETURNING updated_at`
	return r.db.QueryRow(query, user.Name, user.ID).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) List(params models.PaginationParams) ([]models.User, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	listQuery := `SELECT id, email, password_hash, name, role, created_at, updated_at FROM users`

	args := []interface{}{}
	where := ""
	argIdx := 1

	if params.Search != "" {
		where = fmt.Sprintf(` WHERE name ILIKE $%d OR email ILIKE $%d`, argIdx, argIdx+1)
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	if err := r.db.QueryRow(countQuery+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	listQuery += where + fmt.Sprintf(` ORDER BY id ASC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, params.PageSize, params.Offset())

	rows, err := r.db.Query(listQuery, args...)
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
