package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
	"github.com/kaplenko/diplom/internal/runner"
)

type SubmissionService struct {
	repo     *repository.SubmissionRepository
	taskRepo *repository.TaskRepository
	docker   *runner.DockerRunner
	notifier *runner.SubmissionNotifier
}

func NewSubmissionService(
	repo *repository.SubmissionRepository,
	taskRepo *repository.TaskRepository,
	docker *runner.DockerRunner,
	notifier *runner.SubmissionNotifier,
) *SubmissionService {
	return &SubmissionService{
		repo:     repo,
		taskRepo: taskRepo,
		docker:   docker,
		notifier: notifier,
	}
}

// Create saves a Submission with status "pending" and immediately returns it.
// The actual code evaluation runs asynchronously in a background goroutine.
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

	go s.evaluate(sub.ID, req.Code, task.TestCases)

	return sub, nil
}

func (s *SubmissionService) evaluate(submissionID int64, code string, testCases json.RawMessage) {
	ctx := context.Background()

	result, err := s.docker.Run(ctx, code, testCases)
	if err != nil {
		log.Printf("[runner] submission %d: runner error: %v", submissionID, err)
		_ = s.repo.UpdateStatus(submissionID, models.StatusFailed, "internal runner error", 0)
		s.publishUpdate(submissionID)
		return
	}

	status := models.StatusFailed
	if result.Status == "passed" {
		status = models.StatusPassed
	}

	resultJSON, _ := json.Marshal(result)

	if err := s.repo.UpdateStatus(submissionID, status, string(resultJSON), result.Score); err != nil {
		log.Printf("[runner] submission %d: failed to update status: %v", submissionID, err)
	}

	s.publishUpdate(submissionID)
}

func (s *SubmissionService) publishUpdate(submissionID int64) {
	sub, err := s.repo.GetByID(submissionID)
	if err != nil {
		log.Printf("[runner] submission %d: failed to load for notification: %v", submissionID, err)
		return
	}
	if sub != nil {
		s.notifier.Publish(sub)
	}
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
