package models

import (
	"encoding/json"
	"time"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type Task struct {
	ID          int64           `json:"id"`
	LessonID    int64           `json:"lesson_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	InitialCode string          `json:"initial_code"`
	TestCases   json.RawMessage `json:"test_cases" swaggertype:"array,object"`
	Difficulty  Difficulty      `json:"difficulty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"required,min=1"`
	InitialCode string     `json:"initial_code"`
	TestCases   []TestCase `json:"test_cases" binding:"required,min=1,dive"`
	Difficulty  Difficulty `json:"difficulty" binding:"required,oneof=easy medium hard"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description string     `json:"description" binding:"omitempty,min=1"`
	InitialCode *string    `json:"initial_code"`
	TestCases   []TestCase `json:"test_cases" binding:"omitempty,min=1,dive"`
	Difficulty  Difficulty `json:"difficulty" binding:"omitempty,oneof=easy medium hard"`
}
