package models

import "time"

type Progress struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	CourseID    int64      `json:"course_id"`
	LessonID    int64      `json:"lesson_id"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type CourseProgress struct {
	CourseID       int64   `json:"course_id"`
	CourseTitle    string  `json:"course_title"`
	TotalLessons   int     `json:"total_lessons"`
	CompletedCount int     `json:"completed_count"`
	Percentage     float64 `json:"percentage"`
}
