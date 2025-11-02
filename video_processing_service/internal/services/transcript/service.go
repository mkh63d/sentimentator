package transcript

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/audio"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

type Service struct {
	videoService *video.Service
	audioService *audio.Service
	strategies   []Strategy
}

func NewService(videoService *video.Service, audioService *audio.Service, whisperAPIKey string) *Service {
	useMockWhisper := isMockWhisperEnabled()

	var whisperStrategy Strategy
	if useMockWhisper {
		slog.Info("WHISPER_MOCK_MODE enabled - using mock Lorem Ipsum transcripts instead of OpenAI API")
		whisperStrategy = NewMockWhisperStrategy()
	} else {
		slog.Info("Using real Whisper API for transcription")
		whisperStrategy = NewWhisperStrategy(whisperAPIKey, audioService)
	}

	return &Service{
		videoService: videoService,
		audioService: audioService,
		strategies: []Strategy{
			NewCaptionStrategy(videoService),
			whisperStrategy,
		},
	}
}

// isMockWhisperEnabled checks if mock mode is enabled via environment variable.
// Set WHISPER_MOCK_MODE=true to enable mock mode for testing without API calls.
func isMockWhisperEnabled() bool {
	mockMode := os.Getenv("WHISPER_MOCK_MODE")
	return strings.ToLower(strings.TrimSpace(mockMode)) == "true"
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
