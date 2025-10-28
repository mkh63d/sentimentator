package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/queue"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoHandler struct {
	db                database.Service
	videoService      *video.Service
	transcriptService *transcript.Service
	queueService      *queue.Service
}

func NewVideoHandler(
	db database.Service,
	videoService *video.Service,
	transcriptService *transcript.Service,
	queueService *queue.Service,
) *VideoHandler {
	return &VideoHandler{
		db:                db,
		videoService:      videoService,
		transcriptService: transcriptService,
		queueService:      queueService,
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
			Status:      "pending",
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

	if err := h.db.UpdateVideoStatus(ctx, videoRecord.ID, "processing"); err != nil {
		slog.Error("failed to update video status",
			slog.String("error", err.Error()),
			slog.Uint64("db_id", uint64(videoRecord.ID)))
	}

	language := "en" // TODO: get from form input

	hasTranscript, existingTranscript, err := h.db.HasTranscript(ctx, videoRecord.ID, language)
	if err != nil {
		slog.Error("failed to check for existing transcript",
			slog.String("error", err.Error()),
			slog.Uint64("video_id", uint64(videoRecord.ID)))
		component := web.VideoUrlError("Failed to check for existing transcript")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
		return
	}

	var transcriptRecord *models.Transcript

	// Only create new transcript and job if one doesn't exist or if it failed
	if !hasTranscript || existingTranscript.Status == models.StatusFailed {
		transcriptRecord = &models.Transcript{
			VideoID:  videoRecord.ID,
			Language: language,
			Method:   models.MethodWhisper,
			Status:   models.StatusPending,
		}

		if err := h.db.CreateTranscript(ctx, transcriptRecord); err != nil {
			slog.Error("failed to create transcript record",
				slog.String("error", err.Error()),
				slog.Uint64("video_id", uint64(videoRecord.ID)))
			component := web.VideoUrlError("Failed to create transcript record")
			if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
				c.String(http.StatusInternalServerError, renderErr.Error())
			}
			return
		}

		slog.Info("created transcript record",
			slog.Uint64("transcript_id", uint64(transcriptRecord.ID)),
			slog.Uint64("video_id", uint64(videoRecord.ID)))

		jobID := uuid.New().String()
		jobPayload := &queue.TranscriptJobPayload{
			VideoURL:     videoURL,
			Language:     language,
			VideoID:      videoRecord.ID,
			TranscriptID: transcriptRecord.ID,
			JobID:        jobID,
		}

		job, err := h.queueService.CreateTranscriptJob(ctx, jobPayload)
		if err != nil {
			slog.Error("failed to create job",
				slog.String("error", err.Error()),
				slog.Uint64("transcript_id", uint64(transcriptRecord.ID)))
			component := web.VideoUrlError("Failed to queue transcript job")
			if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
				c.String(http.StatusInternalServerError, renderErr.Error())
			}
			return
		}

		slog.Info("created transcript job",
			slog.Uint64("job_id", uint64(job.ID)),
			slog.Uint64("transcript_id", uint64(transcriptRecord.ID)),
			slog.String("tracking_id", jobID))
	} else {
		slog.Info("transcript already exists, skipping job creation",
			slog.Uint64("transcript_id", uint64(existingTranscript.ID)),
			slog.String("status", string(existingTranscript.Status)))
		transcriptRecord = existingTranscript
	}

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
		slog.Uint64("db_id", uint64(videoRecord.ID)),
		slog.Uint64("transcript_id", uint64(transcriptRecord.ID)))
}

func (h *VideoHandler) GetJobStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stats, err := h.db.GetJobStats(ctx)
	if err != nil {
		slog.Error("failed to get job stats", slog.String("error", err.Error()))
		component := web.JobStatsError("Failed to load job statistics")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
		return
	}

	component := web.JobStatsComponent(stats)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render job stats component", slog.String("error", err.Error()))
	}
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
