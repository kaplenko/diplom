package service

import (
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type ProgressService struct {
	repo *repository.ProgressRepository
}

func NewProgressService(repo *repository.ProgressRepository) *ProgressService {
	return &ProgressService{repo: repo}
}

func (s *ProgressService) MarkLessonCompleted(userID, courseID, lessonID int64) error {
	if err := s.repo.MarkCompleted(userID, courseID, lessonID); err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}
	return nil
}

func (s *ProgressService) GetCourseProgress(userID, courseID int64) (*models.CourseProgress, error) {
	cp, err := s.repo.GetCourseProgress(userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("get course progress: %w", err)
	}
	if cp == nil {
		return nil, ErrNotFound
	}
	return cp, nil
}

func (s *ProgressService) GetAllProgress(userID int64) ([]models.CourseProgress, error) {
	return s.repo.GetAllProgress(userID)
}
