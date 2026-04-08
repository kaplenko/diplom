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

// CourseProgress represents overall course completion (lessons completed).
// Percentage is an integer 0-100 suitable for direct use in a progress bar.
type CourseProgress struct {
	CourseID       int64  `json:"course_id"`
	CourseTitle    string `json:"course_title"`
	TotalLessons   int   `json:"total_lessons"`
	CompletedCount int   `json:"completed_count"`
	Percentage     int   `json:"percentage"`
}

// LessonProgress represents how many tasks the user has solved inside a
// single lesson. Percentage is an integer 0-100.
type LessonProgress struct {
	LessonID       int64  `json:"lesson_id"`
	LessonTitle    string `json:"lesson_title"`
	TotalTasks     int    `json:"total_tasks"`
	CompletedTasks int    `json:"completed_tasks"`
	Percentage     int    `json:"percentage"`
}
