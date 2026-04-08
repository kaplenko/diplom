package service

import (
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type CourseService struct {
	repo *repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) Create(req models.CreateCourseRequest, createdBy int64) (*models.Course, error) {
	course := &models.Course{
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(course); err != nil {
		return nil, fmt.Errorf("create course: %w", err)
	}

	return course, nil
}

func (s *CourseService) GetByID(id int64) (*models.Course, error) {
	course, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	if course == nil {
		return nil, ErrNotFound
	}
	return course, nil
}

func (s *CourseService) Update(id int64, req models.UpdateCourseRequest) (*models.Course, error) {
	course, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	if course == nil {
		return nil, ErrNotFound
	}

	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}

	if err := s.repo.Update(course); err != nil {
		return nil, fmt.Errorf("update course: %w", err)
	}

	return course, nil
}

func (s *CourseService) Delete(id int64) error {
	course, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get course: %w", err)
	}
	if course == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}

func (s *CourseService) List(params models.PaginationParams) ([]models.Course, int64, error) {
	return s.repo.List(params)
}
