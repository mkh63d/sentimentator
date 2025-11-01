package audio

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

// Extractor handles audio extraction from video platforms.
type Extractor interface {
	ExtractAudio(ctx context.Context, videoURL string) (string, error)
}

type YtDlpExtractor struct {
	videoService *video.Service
	tempDir      string
}

// NewYtDlpExtractor creates a new yt-dlp based audio extractor
func NewYtDlpExtractor(videoService *video.Service) *YtDlpExtractor {
	return &YtDlpExtractor{
		videoService: videoService,
		tempDir:      os.TempDir(),
	}
}

func (e *YtDlpExtractor) ExtractAudio(ctx context.Context, videoURL string) (string, error) {
	platform, err := e.videoService.DetectPlatform(videoURL)
	if err != nil {
		return "", fmt.Errorf("failed to detect platform: %w", err)
	}

	config := platform.GetDownloadConfig()

	outputPath := filepath.Join(e.tempDir, fmt.Sprintf("temp_audio_%d.mp3", os.Getpid()))

	args := []string{
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
	}

	if config.UserAgent != "" {
		args = append(args, "--user-agent", config.UserAgent)
	}

	if config.Referer != "" {
		args = append(args, "--referer", config.Referer)
	}

	for key, value := range config.Headers {
		args = append(args, "--add-header", fmt.Sprintf("%s:%s", key, value))
	}

	args = append(args, config.ExtraArgs...)

	if config.CookieFile != "" {
		args = append(args, "--cookies", config.CookieFile)
	}

	args = append(args, videoURL)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp command failed for %s: %w, stderr: %s",
			platform.Name(), err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file was not created at %s", outputPath)
	}

	return outputPath, nil
}

func (e *YtDlpExtractor) SetTempDir(dir string) {
	e.tempDir = dir
}
