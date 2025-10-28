package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
)

type WorkerPool struct {
	db                database.Service
	queueService      *Service
	transcriptService *transcript.Service
	workers           int
	pollInterval      time.Duration
	wg                sync.WaitGroup
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewWorkerPool(
	db database.Service,
	transcriptService *transcript.Service,
	workers int,
	pollInterval time.Duration,
) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		db:                db,
		queueService:      NewService(db),
		transcriptService: transcriptService,
		workers:           workers,
		pollInterval:      pollInterval,
		ctx:               ctx,
		cancel:            cancel,
	}
}

func (wp *WorkerPool) Start() {
	slog.Info("starting worker pool", "workers", wp.workers)

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	slog.Info("stopping worker pool")
	wp.cancel()
	wp.wg.Wait()
	slog.Info("worker pool stopped")
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	slog.Info("worker started", "worker_id", id)
	ticker := time.NewTicker(wp.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			slog.Info("worker stopping", "worker_id", id)
			return
		case <-ticker.C:
			wp.processBatch(id)
		}
	}
}

func (wp *WorkerPool) processBatch(workerID int) {
	ctx, cancel := context.WithTimeout(wp.ctx, 5*time.Minute)
	defer cancel()

	jobs, err := wp.queueService.GetPendingJobs(ctx, models.JobTypeTranscript, 1)
	if err != nil {
		slog.Error("failed to fetch pending jobs",
			"worker_id", workerID,
			"error", err.Error())
		return
	}

	if len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		if err := wp.processJob(ctx, workerID, &job); err != nil {
			slog.Error("failed to process job",
				"worker_id", workerID,
				"job_id", job.ID,
				"error", err.Error())
		}
	}
}

func (wp *WorkerPool) processJob(ctx context.Context, workerID int, job *models.Job) error {
	if err := wp.queueService.MarkJobProcessing(ctx, job.ID); err != nil {
		return err
	}

	slog.Info("processing job",
		"worker_id", workerID,
		"job_id", job.ID,
		"type", job.Type,
		"attempts", job.Attempts+1)

	switch job.Type {
	case models.JobTypeTranscript:
		return wp.processTranscriptJob(ctx, job)
	default:
		return wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusFailed, "unsupported job type")
	}
}

func (wp *WorkerPool) processTranscriptJob(ctx context.Context, job *models.Job) error {
	// Parse payload
	var payload TranscriptJobPayload
	if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
		errMsg := "failed to parse job payload: " + err.Error()
		_ = wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusFailed, errMsg)
		return err
	}

	slog.Info("processing transcript job",
		"job_id", job.ID,
		"video_url", payload.VideoURL,
		"transcript_id", payload.TranscriptID)

	transcriptRecord, err := wp.db.GetTranscriptByID(ctx, payload.TranscriptID)
	if err != nil {
		errMsg := "failed to get transcript record: " + err.Error()
		_ = wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusFailed, errMsg)
		return err
	}

	if err := wp.db.UpdateTranscriptStatus(ctx, payload.TranscriptID, models.StatusProcessing); err != nil {
		slog.Warn("failed to update transcript status to processing", "error", err.Error())
	}

	result, err := wp.transcriptService.GetTranscript(ctx, payload.VideoURL, payload.Language)
	if err != nil {
		errMsg := "failed to get transcript: " + err.Error()

		_ = wp.db.UpdateTranscriptStatus(ctx, payload.TranscriptID, models.StatusFailed)

		_ = wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusFailed, errMsg)

		return err
	}

	transcriptRecord.Method = result.Method
	transcriptRecord.Text = result.Text
	transcriptRecord.Segments = result.Segments
	transcriptRecord.Status = models.StatusCompleted
	now := time.Now()
	transcriptRecord.CompletedAt = &now

	if err := wp.db.GetDB().WithContext(ctx).Save(transcriptRecord).Error; err != nil {
		errMsg := "failed to save completed transcript: " + err.Error()
		_ = wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusFailed, errMsg)
		return err
	}

	if err := wp.db.UpdateVideoStatus(ctx, payload.VideoID, "completed"); err != nil {
		slog.Warn("failed to update video status", "error", err.Error())
	}

	if err := wp.queueService.UpdateJobStatus(ctx, job.ID, models.JobStatusCompleted, ""); err != nil {
		slog.Error("failed to mark job as completed", "error", err.Error())
		return err
	}

	slog.Info("transcript job completed successfully",
		"job_id", job.ID,
		"transcript_id", payload.TranscriptID,
		"method", string(result.Method),
		"text_length", len(result.Text))

	return nil
}
