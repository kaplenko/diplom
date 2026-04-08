package service

import (
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(id int64) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *UserService) Update(id int64, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (s *UserService) Delete(id int64) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}

func (s *UserService) List(params models.PaginationParams) ([]models.User, int64, error) {
	return s.repo.List(params)
}
