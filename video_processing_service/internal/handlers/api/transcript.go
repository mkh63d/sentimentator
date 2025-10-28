package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/queue"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TranscriptHandler struct {
	db           database.Service
	videoService *video.Service
	queueService *queue.Service
}

func NewTranscriptHandler(
	db database.Service,
	videoService *video.Service,
	queueService *queue.Service,
) *TranscriptHandler {
	return &TranscriptHandler{
		db:           db,
		videoService: videoService,
		queueService: queueService,
	}
}

func (h *TranscriptHandler) GetTranscript(c *gin.Context) {
	var req models.TranscriptRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.VideoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_url is required"})
		return
	}

	if req.Language == "" {
		req.Language = "en"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	metadata, err := h.videoService.GetVideoMetadata(ctx, req.VideoURL)
	if err != nil {
		slog.Error("failed to get video metadata", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get video metadata: " + err.Error()})
		return
	}

	videoRecord, err := h.db.GetVideoByVideoID(ctx, metadata.VideoID, metadata.Platform)
	if err != nil {
		videoRecord = &models.Video{
			VideoID:     metadata.VideoID,
			Platform:    metadata.Platform,
			OriginalURL: req.VideoURL,
			Title:       metadata.Title,
			Duration:    metadata.Duration,
			AudioURL:    metadata.AudioURL,
			Status:      "pending",
		}

		if err := h.db.SaveVideo(ctx, videoRecord); err != nil {
			slog.Error("failed to save video", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save video"})
			return
		}

		slog.Info("created new video record",
			"video_id", metadata.VideoID,
			"db_id", videoRecord.ID)
	}

	hasTranscript, existingTranscript, err := h.db.HasTranscript(ctx, videoRecord.ID, req.Language)
	if err != nil {
		slog.Error("failed to check for existing transcript", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check for existing transcript"})
		return
	}

	if hasTranscript {
		slog.Info("transcript already exists",
			"transcript_id", existingTranscript.ID,
			"video_id", videoRecord.ID,
			"status", existingTranscript.Status)

		switch existingTranscript.Status {
		case models.StatusCompleted:
			c.JSON(http.StatusOK, models.TranscriptResponse{
				ID:     fmt.Sprintf("%d", existingTranscript.ID),
				Status: existingTranscript.Status,
			})
		case models.StatusProcessing, models.StatusPending:
			c.JSON(http.StatusAccepted, models.TranscriptResponse{
				ID:     fmt.Sprintf("%d", existingTranscript.ID),
				Status: existingTranscript.Status,
			})
		case models.StatusFailed:
			slog.Info("retrying failed transcript",
				"transcript_id", existingTranscript.ID,
				"video_id", videoRecord.ID)
		}

		// Don't create a new job if transcript exists and is not failed
		if existingTranscript.Status != models.StatusFailed {
			return
		}
	}

	transcriptRecord := &models.Transcript{
		VideoID:  videoRecord.ID,
		Language: req.Language,
		Method:   models.MethodWhisper, // TODO: think about default method
		Status:   models.StatusPending,
	}

	if err := h.db.CreateTranscript(ctx, transcriptRecord); err != nil {
		slog.Error("failed to create transcript record", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transcript record"})
		return
	}

	slog.Info("created transcript record",
		"transcript_id", transcriptRecord.ID,
		"video_id", videoRecord.ID)

	jobID := uuid.New().String()
	jobPayload := &queue.TranscriptJobPayload{
		VideoURL:     req.VideoURL,
		Language:     req.Language,
		VideoID:      videoRecord.ID,
		TranscriptID: transcriptRecord.ID,
		JobID:        jobID,
	}

	job, err := h.queueService.CreateTranscriptJob(ctx, jobPayload)
	if err != nil {
		slog.Error("failed to create job", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue transcript job"})
		return
	}

	slog.Info("created transcript job",
		"job_id", job.ID,
		"transcript_id", transcriptRecord.ID,
		"tracking_id", jobID)

	response := models.TranscriptResponse{
		ID:     jobID,
		Status: models.StatusPending,
	}

	c.JSON(http.StatusAccepted, response)
}
