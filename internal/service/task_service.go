package service

import (
	"encoding/json"
	"fmt"

	"github.com/kaplenko/diplom/internal/models"
	"github.com/kaplenko/diplom/internal/repository"
)

type TaskService struct {
	repo       *repository.TaskRepository
	lessonRepo *repository.LessonRepository
}

func NewTaskService(repo *repository.TaskRepository, lessonRepo *repository.LessonRepository) *TaskService {
	return &TaskService{repo: repo, lessonRepo: lessonRepo}
}

func (s *TaskService) Create(lessonID int64, req models.CreateTaskRequest) (*models.Task, error) {
	lesson, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("get lesson: %w", err)
	}
	if lesson == nil {
		return nil, ErrNotFound
	}

	testCasesJSON, err := json.Marshal(req.TestCases)
	if err != nil {
		return nil, fmt.Errorf("marshal test cases: %w", err)
	}

	task := &models.Task{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		InitialCode: req.InitialCode,
		TestCases:   testCasesJSON,
		Difficulty:  req.Difficulty,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (s *TaskService) GetByID(id int64) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if task == nil {
		return nil, ErrNotFound
	}
	return task, nil
}

func (s *TaskService) Update(id int64, req models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if task == nil {
		return nil, ErrNotFound
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.InitialCode != nil {
		task.InitialCode = *req.InitialCode
	}
	if req.Difficulty != "" {
		task.Difficulty = req.Difficulty
	}
	if len(req.TestCases) > 0 {
		testCasesJSON, err := json.Marshal(req.TestCases)
		if err != nil {
			return nil, fmt.Errorf("marshal test cases: %w", err)
		}
		task.TestCases = testCasesJSON
	}

	if err := s.repo.Update(task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

func (s *TaskService) Delete(id int64) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}
	if task == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}

func (s *TaskService) ListByLesson(lessonID int64, params models.PaginationParams) ([]models.Task, int64, error) {
	return s.repo.ListByLesson(lessonID, params)
}
