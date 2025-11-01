package audio

import (
	"context"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

// Service manages audio extraction operations across different platforms.
// This service acts as a facade, providing a simple interface for audio extraction
// while delegating the actual work to platform-aware extractors.
type Service struct {
	extractor Extractor
}

// NewService creates a new audio service with the default YtDlpExtractor.
// The video.Service is injected to enable platform detection and configuration.
func NewService(videoService *video.Service) *Service {
	return &Service{
		extractor: NewYtDlpExtractor(videoService),
	}
}

// ExtractAudio extracts audio from a video URL and returns the path to the audio file.
func (s *Service) ExtractAudio(ctx context.Context, videoURL string) (string, error) {
	return s.extractor.ExtractAudio(ctx, videoURL)
}

func (s *Service) SetExtractor(extractor Extractor) {
	s.extractor = extractor
}
