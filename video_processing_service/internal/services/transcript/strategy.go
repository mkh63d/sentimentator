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
	Text             string
	Segments         []models.Segment
	Method           models.TranscriptMethod
	DetectedLanguage string

	AudioExtractionDuration int64
	TranscriptionDuration   int64
	TotalProcessingDuration int64
	AudioFileSizeBytes      int64
}
