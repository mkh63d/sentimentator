package models

import (
	"time"

	"gorm.io/gorm"
)

type JobType string

const (
	JobTypeTranscript   JobType = "transcript"
	JobTypeAudioExtract JobType = "audio_extract"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Type        JobType        `gorm:"type:varchar(50);not null;index:idx_jobs_type_status,priority:1" json:"type"`
	Payload     string         `gorm:"type:jsonb" json:"payload"`
	Status      JobStatus      `gorm:"type:varchar(50);default:'pending';index:idx_jobs_type_status,priority:2;index:idx_jobs_status_created,priority:1" json:"status"`
	Attempts    int            `gorm:"default:0" json:"attempts"`
	MaxAttempts int            `gorm:"default:3" json:"max_attempts"`
	Error       string         `gorm:"type:text" json:"error,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_jobs_status_created,priority:2" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Job) TableName() string {
	return "jobs"
}
