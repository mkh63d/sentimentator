package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	db database.Service
}

func NewService(db database.Service) *Service {
	return &Service{
		db: db,
	}
}

// CreateTranscriptJob creates a new transcript job in the queue
func (s *Service) CreateTranscriptJob(ctx context.Context, payload *TranscriptJobPayload) (*models.Job, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	job := &models.Job{
		Type:        models.JobTypeTranscript,
		Payload:     string(payloadJSON),
		Status:      models.JobStatusPending,
		Attempts:    0,
		MaxAttempts: 1,
	}

	if err := s.db.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return job, nil

}

func (s *Service) GetPendingJobs(ctx context.Context, jobType models.JobType, limit int) ([]models.Job, error) {
	return s.db.GetPendingJobs(ctx, jobType, limit)
}

// UpdateJobStatus updates the status of a job
func (s *Service) UpdateJobStatus(ctx context.Context, jobID uint, status models.JobStatus, errorMsg string) error {
	switch status {
	case models.JobStatusCompleted:
		return s.db.MarkJobCompleted(ctx, jobID)
	case models.JobStatusFailed:
		return s.db.MarkJobFailed(ctx, jobID, errorMsg)
	default:
		return s.db.MarkJobFailed(ctx, jobID, errorMsg)
	}
}

// MarkJobProcessing marks a job as processing and increments attempts
func (s *Service) MarkJobProcessing(ctx context.Context, jobID uint) error {
	return s.db.GetDB().WithContext(ctx).
		Model(&models.Job{}).
		Where("id = ? AND status = ?", jobID, models.JobStatusPending).
		Updates(map[string]any{
			"status":   models.JobStatusProcessing,
			"attempts": gorm.Expr("attempts + ?", 1),
		}).Error
}
