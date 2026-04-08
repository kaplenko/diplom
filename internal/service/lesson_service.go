package service

import (
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type LessonService struct {
	repo       *repository.LessonRepository
	courseRepo *repository.CourseRepository
}

func NewLessonService(repo *repository.LessonRepository, courseRepo *repository.CourseRepository) *LessonService {
	return &LessonService{repo: repo, courseRepo: courseRepo}
}

func (s *LessonService) Create(courseID int64, req models.CreateLessonRequest) (*models.Lesson, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	if course == nil {
		return nil, ErrNotFound
	}

	lesson := &models.Lesson{
		CourseID:   courseID,
		Title:      req.Title,
		Content:    req.Content,
		OrderIndex: req.OrderIndex,
	}

	if err := s.repo.Create(lesson); err != nil {
		return nil, fmt.Errorf("create lesson: %w", err)
	}

	return lesson, nil
}

func (s *LessonService) GetByID(id int64) (*models.Lesson, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get lesson: %w", err)
	}
	if lesson == nil {
		return nil, ErrNotFound
	}
	return lesson, nil
}

func (s *LessonService) Update(id int64, req models.UpdateLessonRequest) (*models.Lesson, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get lesson: %w", err)
	}
	if lesson == nil {
		return nil, ErrNotFound
	}

	if req.Title != "" {
		lesson.Title = req.Title
	}
	if req.Content != "" {
		lesson.Content = req.Content
	}
	if req.OrderIndex != nil {
		lesson.OrderIndex = *req.OrderIndex
	}

	if err := s.repo.Update(lesson); err != nil {
		return nil, fmt.Errorf("update lesson: %w", err)
	}

	return lesson, nil
}

func (s *LessonService) Delete(id int64) error {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get lesson: %w", err)
	}
	if lesson == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}

func (s *LessonService) ListByCourse(courseID int64, params models.PaginationParams) ([]models.Lesson, int64, error) {
	return s.repo.ListByCourse(courseID, params)
}
