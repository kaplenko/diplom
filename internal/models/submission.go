package models

import "time"

type SubmissionStatus string

const (
	StatusPending SubmissionStatus = "pending"
	StatusPassed  SubmissionStatus = "passed"
	StatusFailed  SubmissionStatus = "failed"
)

type Submission struct {
	ID          int64            `json:"id"`
	TaskID      int64            `json:"task_id"`
	UserID      int64            `json:"user_id"`
	Code        string           `json:"code"`
	Status      SubmissionStatus `json:"status"`
	Result      string           `json:"result"`
	Score       int              `json:"score"`
	SubmittedAt time.Time        `json:"submitted_at"`
}

type CreateSubmissionRequest struct {
	Code string `json:"code" binding:"required,min=1"`
}
