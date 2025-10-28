package transcript

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

type Service struct {
	videoService *video.Service
	strategies   []Strategy
}

func NewService(videoService *video.Service, whisperAPIKey string) *Service {
	return &Service{
		videoService: videoService,
		strategies: []Strategy{
			// NewCaptionStrategy(videoService),
			NewWhisperStrategy(whisperAPIKey),
		},
	}
}

func (s *Service) GetTranscript(ctx context.Context, videoURL, lang string) (*Result, error) {
	metadata, err := s.videoService.GetVideoMetadata(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get video metadata: %w", err)
	}

	var lastErr error
	for _, strategy := range s.strategies {
		slog.Info("Attempting transcript strategy", "strategy", strategy.Name())

		result, err := strategy.GetTranscript(ctx, videoURL, metadata.Platform, metadata.VideoID, lang)
		if err != nil {
			slog.Warn("Transcript strategy failed", "strategy", strategy.Name(), "error", err)
			lastErr = err
			continue
		}

		slog.Info("Transcript strategy succeeded", "strategy", strategy.Name())
		return result, nil
	}

	return nil, fmt.Errorf("all transcript strategies failed: %w", lastErr)
}
