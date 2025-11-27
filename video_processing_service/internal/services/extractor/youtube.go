package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type YouTubePlatform struct {
	client  *http.Client
	tempDir string
}

func NewYouTubePlatform(tempDir string) Platform {
	return &YouTubePlatform{
		client:  &http.Client{},
		tempDir: tempDir}
}

func (ytp *YouTubePlatform) Name() string {
	return "youtube_platform"
}

func (ytp *YouTubePlatform) SupportsURL(videoURL string) bool {
	return strings.Contains(videoURL, "youtube.com") || strings.Contains(videoURL, "youtu.be")
}

func (y *YouTubePlatform) ExtractVideoID(rawVideoURL string) (string, error) {
	if rawVideoURL == "" {
		return "", errors.New("Video URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawVideoURL)
	if err != nil {
		return "", err
	}

	host := strings.TrimPrefix(strings.ToLower(parsedURL.Host), "www.")

	// youtu.be short links
	if host == "youtu.be" {
		videoID := strings.TrimPrefix(parsedURL.Path, "/")
		if videoID != "" {
			return videoID, nil
		}
	}

	// youtube.com URLs
	if strings.Contains(host, "youtube.com") {
		videoID := parsedURL.Query().Get("v")
		if videoID != "" {
			return videoID, nil
		}

		// Check for /embed/ pattern
		if strings.Contains(parsedURL.Path, "/embed/") {
			re := regexp.MustCompile(`/embed/([a-zA-Z0-9_-]+)`)
			matches := re.FindStringSubmatch(parsedURL.Path)
			if len(matches) > 1 {
				return matches[1], nil
			}
		}

		// Check for /v/ pattern
		if strings.Contains(parsedURL.Path, "/v/") {
			re := regexp.MustCompile(`/v/([a-zA-Z0-9_-]+)`)
			matches := re.FindStringSubmatch(parsedURL.Path)
			if len(matches) > 1 {
				return matches[1], nil
			}
		}

		// Check for /watch/ pattern (less common)
		if strings.Contains(parsedURL.Path, "/watch/") {
			re := regexp.MustCompile(`/watch/([a-zA-Z0-9_-]+)`)
			matches := re.FindStringSubmatch(parsedURL.Path)
			if len(matches) > 1 {
				return matches[1], nil
			}
		}
	}

	return "", errors.New("unable to extract YouTube video ID from URL")
}

func (ytp *YouTubePlatform) GetMetadata(ctx context.Context, videoID string) (*VideoMetadata, error) {
	slog.InfoContext(ctx, "Fetching YouTube video metadata", slog.String("video_id", videoID))
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	args := []string{
		"--dump-json",
		"--skip-download",
		"--no-playlist",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"--no-check-certificates",
	}

	cookiesPath := os.Getenv("YOUTUBE_COOKIES_PATH")
	if cookiesPath != "" {
		if _, err := os.Stat(cookiesPath); err == nil {
			slog.InfoContext(ctx, "using YouTube cookies for authentication", slog.String("path", cookiesPath))
			args = append(args, "--cookies", cookiesPath)
		} else {
			slog.WarnContext(ctx, "YOUTUBE_COOKIES_PATH set but file not found",
				slog.String("path", cookiesPath),
				slog.String("error", err.Error()))
		}
	}

	args = append(args, videoURL)
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		slog.ErrorContext(ctx, "yt-dlp metadata fetch failed",
			slog.String("stderr", stderr.String()))
		return nil, fmt.Errorf("failed to fetch metadata via yt-dlp: %w", err)
	}

	var ytdlpData struct {
		ID        string                 `json:"id"`
		Title     string                 `json:"title"`
		Duration  int                    `json:"duration"` // yt-dlp returns duration in seconds
		Thumbnail string                 `json:"thumbnail"`
		Subtitles map[string]interface{} `json:"subtitles"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &ytdlpData); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}

	hasCaptions := len(ytdlpData.Subtitles) > 0

	thumbnailURL := ytdlpData.Thumbnail
	if strings.Trim(thumbnailURL, "") == "" {
		slog.WarnContext(ctx, "thumbnail URL is empty in yt-dlp metadata", slog.String("video_id", videoID))
		thumbnailURL = fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", videoID)
	}

	return &VideoMetadata{
		VideoID:                   videoID,
		Platform:                  ytp.Name(),
		Title:                     ytdlpData.Title,
		Duration:                  ytdlpData.Duration,
		ThumbnailURL:              thumbnailURL,
		HasCaptions:               hasCaptions,
		AvailableCaptionLanguages: make([]string, 0),
		OriginalURL:               videoURL,
	}, nil
}

func (ytp *YouTubePlatform) GetCaptions(ctx context.Context, videoID, lang string) ([]Segment, error) {
	return make([]Segment, 0), nil
}

func (ytp *YouTubePlatform) ExtractAudio(ctx context.Context, videoID string) (*AudioResult, error) {
	outputPath := filepath.Join(ytp.tempDir, fmt.Sprintf("temp_youtube_audio_%s_%d.mp3", videoID, time.Now().UnixNano()))

	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	args := []string{
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"--referer", "https://www.youtube.com/",
		"--add-header", "Accept-Language: en-US,en;q=0.9",
		"--no-check-certificates",
		"--extractor-args", "youtube:player_client=default,web",
		"--retries", "2",
		"--fragment-retries", "3",
		"--no-cache-dir",
	}

	// Check for YouTube cookies for authentication because when service is running in production, it for sure will be marked as a bot
	cookiesPath := os.Getenv("YOUTUBE_COOKIES_PATH")
	if cookiesPath != "" {
		if _, err := os.Stat(cookiesPath); err == nil {
			slog.InfoContext(ctx, "using YouTube cookies for authentication", slog.String("path", cookiesPath))
			args = append(args, "--cookies", cookiesPath)
		} else {
			slog.WarnContext(ctx, "YOUTUBE_COOKIES_PATH set but file not found",
				slog.String("path", cookiesPath),
				slog.String("error", err.Error()))
		}
	}

	args = append(args, videoURL)
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed for %s: %w, stderr: %s",
			ytp.Name(), err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file was not created at %s", outputPath)
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat audio file: %w", err)
	}

	return &AudioResult{
		FilePath: outputPath,
		Format:   "mp3",
		Duration: -1,
		FileSize: fileInfo.Size(),
		Metadata: nil,
	}, nil
}
