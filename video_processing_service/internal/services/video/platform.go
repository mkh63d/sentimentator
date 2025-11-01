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

// DownloadConfig contains platform-specific configuration for audio/video download tools.
// This struct enables each platform to specify its own requirements for downloads.
type DownloadConfig struct {
	Headers    map[string]string
	Referer    string
	ExtraArgs  []string
	UserAgent  string
	CookieFile string
}

type Platform interface {
	// Name returns the platform identifier (e.g., "youtube", "instagram")
	Name() string

	// SupportsURL checks if this platform can handle the given URL
	SupportsURL(url string) bool

	// ExtractVideoID extracts the unique video identifier from a URL
	ExtractVideoID(url string) (string, error)

	// GetMetadata fetches video metadata (title, duration, caption availability, etc.)
	GetMetadata(ctx context.Context, videoID string) (*models.VideoMetadata, error)

	// GetCaptions retrieves native captions/subtitles if available
	// Returns ErrNoCaptionsAvailable if the platform doesn't support captions
	GetCaptions(ctx context.Context, videoID string, lang string) ([]models.Segment, error)

	// GetAudioURL returns the direct URL or formatted URL for audio extraction
	// TODO: This may be deprecated in favor of GetDownloadConfig
	GetAudioURL(ctx context.Context, videoID string) (string, error)

	// GetDownloadConfig returns platform-specific configuration for audio download.
	GetDownloadConfig() *DownloadConfig
}
