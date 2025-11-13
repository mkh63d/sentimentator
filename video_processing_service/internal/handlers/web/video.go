package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/messaging"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	db                database.Service
	extractorService  extractor.Service
	transcriptService *transcript.Service
	producer          *messaging.Producer
}

func NewVideoHandler(
	db database.Service,
	extractorService extractor.Service,
	transcriptService *transcript.Service,
	producer *messaging.Producer,
) *VideoHandler {
	return &VideoHandler{
		db:                db,
		extractorService:  extractorService,
		transcriptService: transcriptService,
		producer:          producer,
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
