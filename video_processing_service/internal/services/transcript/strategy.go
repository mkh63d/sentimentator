package transcript

import (
	"context"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type Strategy interface {
	Name() string
	GetTranscript(ctx context.Context, videoURL, platform, videoID, lang string) (*Result, error)
}

type Result struct {
	Text     string
	Segments []models.Segment
	Method   models.TranscriptMethod
}
