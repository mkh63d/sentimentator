package web

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
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

// calculateMaxPage calculates the maximum page number based on total items and items per page.
// Returns at least 1 to ensure there's always a valid page number.
func calculateMaxPage(total int64, itemsPerPage int) int {
	maxPage := (int(total) + itemsPerPage - 1) / itemsPerPage
	if maxPage < 1 {
		return 1
	}
	return maxPage
}

func (h *VideoHandler) ListVideos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	const itemsPerPage = 10
	pageParam := c.DefaultQuery("page", "1")
	currentParam := c.DefaultQuery("current", "")

	currentPage := 1
	if pageNum, err := strconv.Atoi(pageParam); err == nil && pageNum > 0 {
		currentPage = pageNum
	} else if cur, err := strconv.Atoi(currentParam); err == nil && cur > 0 {
		currentPage = cur
	}

	renderListError := func() {
		slog.Error("failed to list videos")
		component := web.VideoListError("Failed to load videos")
		if renderErr := component.Render(c.Request.Context(), c.Writer); renderErr != nil {
			c.String(http.StatusInternalServerError, renderErr.Error())
		}
	}

	if pageParam == "last" {
		_, total, err := h.db.ListVideos(ctx, itemsPerPage, 0)
		if err != nil {
			slog.Error("failed to list videos", slog.String("error", err.Error()))
			renderListError()
			return
		}
		currentPage = calculateMaxPage(total, itemsPerPage)
	}

	switch pageParam {
	case "next":
		currentPage++
	case "prev":
		if currentPage > 1 {
			currentPage--
		}
	case "first":
		currentPage = 1
	}

	if currentPage < 1 {
		currentPage = 1
	}

	offset := (currentPage - 1) * itemsPerPage

	videos, total, err := h.db.ListVideos(ctx, itemsPerPage, offset)
	if err != nil {
		slog.Error("failed to list videos", slog.String("error", err.Error()))
		renderListError()
		return
	}

	maxPage := calculateMaxPage(total, itemsPerPage)
	if currentPage > maxPage {
		currentPage = maxPage
		offset = (currentPage - 1) * itemsPerPage
		videos, total, err = h.db.ListVideos(ctx, itemsPerPage, offset)
		if err != nil {
			slog.Error("failed to list videos", slog.String("error", err.Error()))
			renderListError()
			return
		}
	}

	c.Set("currentPage", currentPage)
	c.Set("itemsPerPage", itemsPerPage)
	c.Set("maxPage", maxPage)

	component := web.VideoListComponent(videos, int(total), currentPage, itemsPerPage)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render video list component", slog.String("error", err.Error()))
	}
}
