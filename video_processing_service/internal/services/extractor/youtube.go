package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Cookie", "CONSENT=YES+")

	resp, err := ytp.client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch YouTube page", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "YouTube returned non-200 status",
			slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "fetched YouTube page", slog.Int("body_size", len(body)))

	patterns := []string{
		`(?s)var ytInitialPlayerResponse\s*=\s*({.+?});`,
		`(?s)ytInitialPlayerResponse\s*=\s*({.+?});`,
		`(?s)"playerResponse":"({.+?})"`,
	}

	var playerResponseJSON []byte
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindSubmatch(body)
		if len(matches) >= 2 {
			playerResponseJSON = matches[1]
			slog.InfoContext(ctx, "found player response with pattern", slog.String("pattern", pattern))
			break
		}
	}

	if playerResponseJSON == nil {
		slog.ErrorContext(ctx, "could not find ytInitialPlayerResponse in page")
		slog.DebugContext(ctx, "page preview", slog.String("preview", string(body[:min(500, len(body))])))
		return nil, errors.New("could not extract player response from YouTube page")
	}

	var data struct {
		VideoDetails struct {
			VideoID       string `json:"videoId"`
			Title         string `json:"title"`
			LengthSeconds string `json:"lengthSeconds"`
			Thumbnail     struct {
				Thumbnails []struct {
					URL string `json:"url"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
		} `json:"videoDetails"`
		Captions struct {
			PlayerCaptionsTracklistRenderer struct {
				CaptionTracks []struct {
					BaseURL string `json:"baseUrl"`
				} `json:"captionTracks"`
			} `json:"playerCaptionsTracklistRenderer"`
		} `json:"captions"`
		Microformat struct {
			PlayerMicroformatRenderer struct {
				Title struct {
					SimpleText string `json:"simpleText"`
				} `json:"title"`
			} `json:"playerMicroformatRenderer"`
		} `json:"microformat"`
	}

	if err := json.Unmarshal(playerResponseJSON, &data); err != nil {
		slog.ErrorContext(ctx, "failed to parse player response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to parse video data: %w", err)
	}

	hasCaptions := len(data.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks) > 0

	title := strings.TrimSpace(data.VideoDetails.Title)
	if title == "" {
		if t := strings.TrimSpace(data.Microformat.PlayerMicroformatRenderer.Title.SimpleText); t != "" {
			title = t
		} else {
			// Fallback: parse og:title from HTML
			reOg := regexp.MustCompile(`(?i)<meta\s+property=["']og:title["']\s+content=["']([^"']+)["']`)
			if m := reOg.FindSubmatch(body); len(m) == 2 {
				title = strings.TrimSpace(string(m[1]))
			}
		}
	}

	slog.InfoContext(ctx, "successfully parsed metadata",
		slog.String("video_id", videoID),
		slog.String("title", title),
		slog.Bool("has_captions", hasCaptions))

	return &VideoMetadata{
		VideoID:                   videoID,
		Platform:                  ytp.Name(),
		Title:                     title,
		Duration:                  -1, //TODO: scrape for Duration
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
