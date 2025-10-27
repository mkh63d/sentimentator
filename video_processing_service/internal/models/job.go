package models

import "time"

type JobTypes string

const (
	JobTypeTranscript JobTypes = "transcript"
)

type Job struct {
	ID          string    `json:"id" db:"id"`
	Type        JobTypes  `json:"type" db:"type"`
	Payload     string    `json:"payload" db:"payload"`
	Status      string    `json:"status" db:"status"`
	Attempts    int       `json:"attempts" db:"attempts"`
	MaxAttempts int       `json:"max_attempts" db:"max_attempts"`
	Error       string    `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
