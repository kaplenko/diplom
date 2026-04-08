package models

import "time"

type Course struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"required,min=1,max=5000"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,min=1,max=5000"`
}
