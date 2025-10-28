package video

import (
	"context"
	"errors"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrVideoNotFound       = errors.New("video not found")
	ErrPrivateVideo        = errors.New("video is private")
)

type Platform interface {
	Name() string
	SupportsURL(url string) bool
	ExtractVideoID(url string) (string, error)
	GetMetadata(ctx context.Context, videoID string) (*models.VideoMetadata, error)
	GetCaptions(ctx context.Context, videoID string, lang string) ([]models.Segment, error)
	GetAudioURL(ctx context.Context, videoID string) (string, error)
}
