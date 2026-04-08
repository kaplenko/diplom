package service

import (
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type SubmissionService struct {
	repo     *repository.SubmissionRepository
	taskRepo *repository.TaskRepository
}

func NewSubmissionService(repo *repository.SubmissionRepository, taskRepo *repository.TaskRepository) *SubmissionService {
	return &SubmissionService{repo: repo, taskRepo: taskRepo}
}

func (s *SubmissionService) Create(taskID, userID int64, req models.CreateSubmissionRequest) (*models.Submission, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if task == nil {
		return nil, ErrNotFound
	}

	sub := &models.Submission{
		TaskID: taskID,
		UserID: userID,
		Code:   req.Code,
		Status: models.StatusPending,
	}

	if err := s.repo.Create(sub); err != nil {
		return nil, fmt.Errorf("create submission: %w", err)
	}

	return sub, nil
}

func (s *SubmissionService) GetByID(id int64) (*models.Submission, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get submission: %w", err)
	}
	if sub == nil {
		return nil, ErrNotFound
	}
	return sub, nil
}

func (s *SubmissionService) ListByTask(taskID, userID int64, params models.PaginationParams) ([]models.Submission, int64, error) {
	return s.repo.ListByTask(taskID, userID, params)
}

func (s *SubmissionService) ListAll(params models.PaginationParams) ([]models.Submission, int64, error) {
	return s.repo.ListAll(params)
}
