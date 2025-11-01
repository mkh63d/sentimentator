package web

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/messaging"
	"github.com/gin-gonic/gin"
)

type DebugHandler struct {
	producer *messaging.Producer
}

func NewDebugHandler(producer *messaging.Producer) *DebugHandler {
	return &DebugHandler{
		producer: producer,
	}
}

// ServeDebugPage serves the debug events page
func (h *DebugHandler) ServeDebugPage(c *gin.Context) {
	component := web.DebugEventsPage()
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render debug events page", slog.String("error", err.Error()))
	}
}

// PublishVideoReceived handles publishing a video received event

func (h *DebugHandler) PublishVideoReceived(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {

		slog.Error("failed to parse form", slog.String("error", err.Error()))

		h.renderError(c, "video.received", "Failed to parse form data")

		return

	}

	event := messaging.VideoTranscriptRequestedEvent{
		AnalysisID: c.Request.FormValue("analysis_id"),
		URL:        c.Request.FormValue("url"),

		Timestamp: time.Now(),
	}

	if err := h.producer.PublishVideoReceived(c.Request.Context(), event); err != nil {

		slog.Error("failed to publish transcript requested event",
			slog.String("error", err.Error()),

			slog.String("analysis_id", event.AnalysisID))
		h.renderError(c, "video.received", err.Error())

		return

	}

	slog.Info("published transcript requested event via debug UI",

		slog.String("analysis_id", event.AnalysisID),
		slog.String("url", event.URL))

	h.renderSuccess(c, "video.received")

}

// PublishTranscriptReady handles publishing a transcript ready event

func (h *DebugHandler) PublishTranscriptReady(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {

		slog.Error("failed to parse form", slog.String("error", err.Error()))

		h.renderError(c, "transcript.ready", "Failed to parse form data")

		return

	}

	durationStr := c.Request.FormValue("duration")

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		duration = 0
	}

	text := c.Request.FormValue("text")

	event := messaging.VideoTranscriptCompletedEvent{
		AnalysisID: c.Request.FormValue("analysis_id"),
		VideoID:    c.Request.FormValue("video_id"),

		TranscriptID: c.Request.FormValue("transcript_id"),

		Text: text,

		Language: c.Request.FormValue("language"),

		Duration: duration,

		CreatedAt: time.Now(),
	}

	if err := h.producer.PublishTranscriptReady(c.Request.Context(), event); err != nil {

		slog.Error("failed to publish transcript completed event",
			slog.String("error", err.Error()),

			slog.String("video_id", event.VideoID))

		h.renderError(c, "transcript.ready", err.Error())

		return

	}

	slog.Info("published transcript completed event via debug UI",

		slog.String("video_id", event.VideoID),

		slog.String("transcript_id", event.TranscriptID))

	h.renderSuccess(c, "transcript.ready")

}

// PublishProcessingFailed handles publishing a processing failed event

func (h *DebugHandler) PublishProcessingFailed(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {

		slog.Error("failed to parse form", slog.String("error", err.Error()))

		h.renderError(c, "processing.failed", "Failed to parse form data")

		return

	}

	attemptsStr := c.Request.FormValue("attempts")

	attempts, err := strconv.Atoi(attemptsStr)

	if err != nil {

		attempts = 1

	}

	event := messaging.VideoTranscriptFailedEvent{
		AnalysisID: c.Request.FormValue("analysis_id"),
		VideoID:    c.Request.FormValue("video_id"),
		Stage:      c.Request.FormValue("stage"),

		Error: c.Request.FormValue("error"),

		Attempts: attempts,

		FailedAt: time.Now(),
	}

	if err := h.producer.PublishProcessingFailed(c.Request.Context(), event); err != nil {

		slog.Error("failed to publish processing failed event",

			slog.String("error", err.Error()),

			slog.String("video_id", event.VideoID))

		h.renderError(c, "processing.failed", err.Error())

		return

	}

	slog.Info("published processing failed event via debug UI",

		slog.String("video_id", event.VideoID),

		slog.String("stage", event.Stage))

	h.renderSuccess(c, "processing.failed")

}

func (h *DebugHandler) renderSuccess(c *gin.Context, eventType string) {
	component := web.EventPublishedSuccess(eventType)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render success component", slog.String("error", err.Error()))
	}
}

func (h *DebugHandler) renderError(c *gin.Context, eventType string, errorMsg string) {
	component := web.EventPublishedError(eventType, errorMsg)
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		slog.Error("failed to render error component", slog.String("error", err.Error()))
	}
}
