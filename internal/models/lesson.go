package models

import "time"

type Lesson struct {
	ID         int64     `json:"id"`
	CourseID   int64     `json:"course_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateLessonRequest struct {
	Title      string `json:"title" binding:"required,min=1,max=200"`
	Content    string `json:"content" binding:"required,min=1"`
	OrderIndex int    `json:"order_index" binding:"min=0"`
}

type UpdateLessonRequest struct {
	Title      string `json:"title" binding:"omitempty,min=1,max=200"`
	Content    string `json:"content" binding:"omitempty,min=1"`
	OrderIndex *int   `json:"order_index" binding:"omitempty,min=0"`
}
