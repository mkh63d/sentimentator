package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TranscriptHandler struct {
	service *transcript.Service
}

func NewTranscriptHandler(service *transcript.Service) *TranscriptHandler {
	return &TranscriptHandler{
		service: service,
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

	jobID := uuid.New().String()

	response := models.TranscriptResponse{
		ID:     jobID,
		Status: models.StatusPending,
	}

	c.JSON(http.StatusAccepted, response)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)

	go func() {
		defer cancel()
		h.processTranscript(ctx, jobID, req)
	}()

}

func (h *TranscriptHandler) processTranscript(
	ctx context.Context,
	jobID string,
	req models.TranscriptRequest,
) {
	slog.Info("Starting transcript processing", "jobID", jobID, "videoURL", req.VideoURL)

	transcript, err := h.service.GetTranscript(ctx, req.VideoURL, req.Language)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get transcript",
			slog.String("error", err.Error()),
			slog.String("job_id", jobID))
		return
	}

	// TODO rm this
	// Ensure output directory exists
	outDir := "transcripts"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		slog.ErrorContext(ctx, "failed to create transcripts directory",
			slog.String("error", err.Error()),
			slog.String("job_id", jobID))
		return
	}
	outPath := fmt.Sprintf("%s/%s.txt", outDir, jobID)
	if err := os.WriteFile(outPath, []byte(transcript.Text), 0o644); err != nil {
		slog.ErrorContext(ctx, "failed to write transcript file",
			slog.String("error", err.Error()),
			slog.String("job_id", jobID),
			slog.String("path", outPath))
		return
	}

	slog.InfoContext(ctx, "transcript completed successfully",
		slog.String("method", string(transcript.Method)),
		slog.String("job_id", jobID),
		slog.String("file_path", outPath),
		slog.Int("text_length", len(transcript.Text)))
}

func (h *TranscriptHandler) GetStatus(c *gin.Context) {
	jobID := c.Query("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	// TODO fetch actual status from storage
	response := models.TranscriptResponse{
		ID:     jobID,
		Status: models.StatusCompleted,
	}

	c.JSON(http.StatusOK, response)
}
