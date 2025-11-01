package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/messaging"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	db                database.Service
	videoService      *video.Service
	transcriptService *transcript.Service
	producer          *messaging.Producer
}

func NewVideoHandler(
	db database.Service,
	videoService *video.Service,
	transcriptService *transcript.Service,
	producer *messaging.Producer,
) *VideoHandler {
	return &VideoHandler{
		db:                db,
		videoService:      videoService,
		transcriptService: transcriptService,
		producer:          producer,
	}
}

func (h *VideoHandler) ProcessVideo(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		slog.Error("failed to parse form", slog.String("error", err.Error()))
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	videoURL := c.Request.FormValue("yt_url_input")
	if videoURL == "" {
		slog.Warn("video URL is required")
		c.String(http.StatusBadRequest, "Video URL is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	metadata, err := h.videoService.GetVideoMetadata(ctx, videoURL)
	if err != nil {
		slog.Error("failed to get video metadata",
			slog.String("error", err.Error()),
			slog.String("url", videoURL))
		component := web.VideoUrlError(err.Error())
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
			slog.Error("failed to render error component", slog.String("error", renderErr.Error()))
		}
		return
	}

	existingVideo, err := h.db.GetVideoByVideoID(ctx, metadata.VideoID, metadata.Platform)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed to check existing video",
			slog.String("error", err.Error()),
			slog.String("video_id", metadata.VideoID))
		component := web.VideoUrlError("Failed to check existing video")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
		return
	}

	var videoRecord *models.Video

	if existingVideo != nil {
		slog.Info("video already exists in database",
			slog.String("video_id", metadata.VideoID),
			slog.Uint64("db_id", uint64(existingVideo.ID)))
		videoRecord = existingVideo
	} else {
		videoRecord = &models.Video{
			VideoID:     metadata.VideoID,
			Platform:    metadata.Platform,
			OriginalURL: videoURL,
			Title:       metadata.Title,
			Duration:    metadata.Duration,
			AudioURL:    metadata.AudioURL,
		}

		if err := h.db.SaveVideo(ctx, videoRecord); err != nil {
			slog.Error("failed to save video to database",
				slog.String("error", err.Error()),
				slog.String("video_id", metadata.VideoID))
			component := web.VideoUrlError("Failed to save video to database")
			if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
				c.String(http.StatusInternalServerError, renderErr.Error())
			}
			return
		}

		slog.Info("video saved to database",
			slog.String("video_id", metadata.VideoID),
			slog.Uint64("db_id", uint64(videoRecord.ID)))
	}

	event := messaging.VideoTranscriptRequestedEvent{
		AnalysisID: metadata.VideoID,
		URL:        videoURL,

		Timestamp: time.Now(),
	}

	if err := h.producer.PublishVideoReceived(ctx, event); err != nil {
		slog.Error("failed to publish video received event",
			slog.String("error", err.Error()),
			slog.String("video_id", metadata.VideoID))
		component := web.VideoUrlError("Failed to queue video for processing")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
		return
	}

	slog.Info("published video received event to RabbitMQ",
		slog.String("video_id", metadata.VideoID),
		slog.Uint64("db_id", uint64(videoRecord.ID)))

	component := web.VideoUrlPost(
		metadata.VideoID,
		metadata.Platform,
		videoURL,
		int(videoRecord.ID),
	)

	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render success component", slog.String("error", err.Error()))
		return
	}

	slog.Info("video processing request completed successfully",
		slog.String("video_id", metadata.VideoID),
		slog.Uint64("db_id", uint64(videoRecord.ID)))
}

func (h *VideoHandler) ListVideos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	limit := 20
	offset := 0

	videos, total, err := h.db.ListVideos(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to list videos", slog.String("error", err.Error()))
		component := web.VideoListError("Failed to load videos")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
		return
	}

	component := web.VideoListComponent(videos, int(total))
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render video list component", slog.String("error", err.Error()))
	}
}
