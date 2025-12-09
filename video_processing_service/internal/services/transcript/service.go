package transcript

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
)

type Service struct {
	extractorService extractor.Service
	strategies       []Strategy
}

func NewService(extractorService extractor.Service) *Service {
	useMockWhisper := isMockWhisperEnabled()

	var whisperStrategy Strategy
	if useMockWhisper {
		slog.Info("WHISPER_MOCK_MODE enabled - using mock Lorem Ipsum transcripts instead of OpenAI API")
		whisperStrategy = NewMockWhisperStrategy()
	} else {
		slog.Info("Using real Whisper API for transcription")
		whisperStrategy = NewWhisperStrategy(extractorService)
	}

	return &Service{
		extractorService: extractorService,
		strategies: []Strategy{
			NewCaptionStrategy(extractorService),
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
	return s.GetTranscriptWithAPIKey(ctx, videoURL, lang, "")
}

// GetTranscriptWithAPIKey attempts to get a transcript using user-provided API key for Whisper
func (s *Service) GetTranscriptWithAPIKey(ctx context.Context, videoURL, lang, userAPIKey string) (*Result, error) {
	metadata, err := s.extractorService.GetMetadata(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get video metadata: %w", err)
	}

	var lastErr error
	for _, strategy := range s.strategies {
		slog.Info("Attempting transcript strategy", "strategy", strategy.Name())

		var result *Result
		// If strategy is Whisper and we have a user API key, use the extended method
		if strategy.Name() == "whisper" && userAPIKey != "" {
			if whisperStrat, ok := strategy.(*WhisperStrategy); ok {
				result, err = whisperStrat.GetTranscriptWithAPIKey(ctx, videoURL, metadata.Platform, metadata.VideoID, lang, userAPIKey)
			} else {
				result, err = strategy.GetTranscript(ctx, videoURL, metadata.Platform, metadata.VideoID, lang)
			}
		} else {
			result, err = strategy.GetTranscript(ctx, videoURL, metadata.Platform, metadata.VideoID, lang)
		}

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
