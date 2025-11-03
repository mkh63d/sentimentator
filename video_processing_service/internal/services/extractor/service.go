package extractor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Service interface {
	ExtractAudio(ctx context.Context, videoURL string) (*AudioResult, error)
	GetMetadata(ctx context.Context, videoURL string) (*VideoMetadata, error)
	GetCaptions(ctx context.Context, videoURL string, lang string) ([]Segment, error)
	DetectPlatform(videoURL string) (string, error)
	SupportedPlatforms() []string
}

type Platform interface {
	// Name returns the platform identifier (e.g., "youtube", "instagram")
	Name() string
	// SupportsURL checks if this platform can handle the given URL
	SupportsURL(url string) bool
	// ExtractVideoID extracts the unique video identifier from a URL
	ExtractVideoID(url string) (string, error)
	GetMetadata(ctx context.Context, videoID string) (*VideoMetadata, error)
	// GetCaptions retrieves native captions/subtitles if available
	// Returns ErrNoCaptionsAvailable if the platform doesn't support captions
	GetCaptions(ctx context.Context, videoID string, lang string) ([]Segment, error)
	ExtractAudio(ctx context.Context, videoID string) (*AudioResult, error)
}

type AudioResult struct {
	FilePath string
	Format   string
	Duration int
	FileSize int64
	Metadata *VideoMetadata
}

type VideoMetadata struct {
	VideoID                   string
	Platform                  string
	Title                     string
	Duration                  int
	HasCaptions               bool
	AvailableCaptionLanguages []string
	OriginalURL               string
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrVideoNotFound       = errors.New("video not found")
	ErrPrivateVideo        = errors.New("video is private or requires authentication")
	ErrNoCaptionsAvailable = errors.New("no captions available for this video")
	ErrExtractionFailed    = errors.New("audio extraction failed")
	ErrInvalidURL          = errors.New("invalid video URL")
)

type extractorService struct {
	platforms map[string]Platform
	tempDir   string
}

func NewService(tempDir string) Service {
	td := tempDir
	if strings.Trim(tempDir, "") == "" {
		td = os.TempDir()
	}

	s := &extractorService{
		platforms: make(map[string]Platform),
		tempDir:   td,
	}

	s.RegisterPlatform(NewYouTubePlatform(td))

	return s
}

func (s *extractorService) RegisterPlatform(platform Platform) {
	s.platforms[platform.Name()] = platform
}

func (s *extractorService) ExtractAudio(ctx context.Context, videoURL string) (*AudioResult, error) {
	platform, err := s.detectPlatform(videoURL)
	if err != nil {
		return nil, err
	}

	videoID, err := platform.ExtractVideoID(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %w", err)
	}

	result, err := platform.ExtractAudio(ctx, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to extract audio from %s: %w", platform.Name(), err)
	}

	return result, nil
}

func (s *extractorService) GetMetadata(ctx context.Context, videoURL string) (*VideoMetadata, error) {
	platform, err := s.detectPlatform(videoURL)
	if err != nil {
		return nil, err
	}

	videoID, err := platform.ExtractVideoID(videoURL)
	if err != nil {
		return nil, err
	}

	return platform.GetMetadata(ctx, videoID)
}

func (s *extractorService) GetCaptions(ctx context.Context, videoURL string, lang string) ([]Segment, error) {
	platform, err := s.detectPlatform(videoURL)
	if err != nil {
		return nil, err
	}

	videoID, err := platform.ExtractVideoID(videoURL)
	if err != nil {
		return nil, err
	}

	return platform.GetCaptions(ctx, videoID, lang)
}

func (s *extractorService) detectPlatform(videoURL string) (Platform, error) {
	for _, platform := range s.platforms {
		if platform.SupportsURL(videoURL) {
			return platform, nil
		}
	}
	return nil, ErrUnsupportedPlatform
}

func (s *extractorService) DetectPlatform(videoURL string) (string, error) {
	platform, err := s.detectPlatform(videoURL)
	if err != nil {
		return "", err
	}
	return platform.Name(), nil
}

func (s *extractorService) SupportedPlatforms() []string {
	platforms := make([]string, 0, len(s.platforms))
	for name := range s.platforms {
		platforms = append(platforms, name)
	}
	return platforms
}
