package video

import (
	"context"
	"fmt"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type Service struct {
	platforms []Platform
}

func NewService() *Service {
	return &Service{
		platforms: []Platform{
			NewYouTubePlatform(),
		},
	}
}

func (s *Service) DetectPlatform(rawVideoURL string) (Platform, error) {
	for _, platform := range s.platforms {
		if platform.SupportsURL(rawVideoURL) {
			return platform, nil
		}
	}

	return nil, ErrUnsupportedPlatform
}

func (s *Service) GetVideoMetadata(ctx context.Context, rawVideoURL string) (*models.VideoMetadata, error) {
	platform, err := s.DetectPlatform(rawVideoURL)
	if err != nil {
		return nil, err
	}

	videoID, err := platform.ExtractVideoID(rawVideoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %w", err)
	}

	metadata, err := platform.GetMetadata(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video metadata: %w", err)
	}

	return metadata, nil
}
