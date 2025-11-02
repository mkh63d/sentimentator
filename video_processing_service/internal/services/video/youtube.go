package video

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type YouTubePlatform struct {
	client *http.Client
}

func NewYouTubePlatform() *YouTubePlatform {
	return &YouTubePlatform{
		client: &http.Client{},
	}
}

func (y *YouTubePlatform) Name() string {
	return "youtube"
}

func (y *YouTubePlatform) SupportsURL(videoURL string) bool {
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

func (y *YouTubePlatform) GetMetadata(ctx context.Context, videoID string) (*models.VideoMetadata, error) {
	slog.InfoContext(ctx, "Fetching YouTube video metadata", slog.String("video_id", videoID))
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := y.client.Do(req)
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
		`var ytInitialPlayerResponse = ({.+?});`,
		`ytInitialPlayerResponse\s*=\s*({.+?});`,
		`"playerResponse":"({.+?})"`,
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
	}

	if err := json.Unmarshal(playerResponseJSON, &data); err != nil {
		slog.ErrorContext(ctx, "failed to parse player response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to parse video data: %w", err)
	}

	hasCaptions := len(data.Captions.PlayerCaptionsTracklistRenderer.CaptionTracks) > 0

	slog.InfoContext(ctx, "successfully parsed metadata",
		slog.String("title", data.VideoDetails.Title),
		slog.Bool("has_captions", hasCaptions))

	return &models.VideoMetadata{
		VideoID:     videoID,
		Platform:    "youtube",
		Title:       data.VideoDetails.Title,
		HasCaptions: hasCaptions,
	}, nil
}

func (y *YouTubePlatform) GetCaptions(ctx context.Context, videoID string, lang string) ([]models.Segment, error) {
	if videoID == "" {
		return nil, errors.New("videoID cannot be empty")
	}

	if lang == "" {
		lang = "en"
	}

	captionURL, err := y.getCaptionURL(ctx, videoID, lang)
	if err != nil {
		return nil, err
	}

	return y.fetchCaptions(ctx, captionURL)
}

func (y *YouTubePlatform) GetAudioURL(
	ctx context.Context,
	videoID string,
) (string, error) {
	// TODO: implement audio extraction for YouTube
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil
}

// TODO: test this way of getting captions, because this part was vibecoded

func (y *YouTubePlatform) getCaptionURL(
	ctx context.Context,
	videoID, lang string,
) (string, error) {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := y.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("no captions available")
	}

	var tracks []struct {
		BaseURL      string `json:"baseUrl"`
		LanguageCode string `json:"languageCode"`
	}

	if err := json.Unmarshal(matches[1], &tracks); err != nil {
		return "", err
	}

	for _, track := range tracks {
		if track.LanguageCode == lang {
			return track.BaseURL, nil
		}
	}

	if len(tracks) > 0 {
		return tracks[0].BaseURL, nil
	}

	return "", fmt.Errorf("no captions found")
}

func (y *YouTubePlatform) fetchCaptions(
	ctx context.Context,
	captionURL string,
) ([]models.Segment, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", captionURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseYouTubeCaptions(body)
}

func parseYouTubeCaptions(data []byte) ([]models.Segment, error) {
	// Parse YouTube's timed text format (simplified)
	var segments []models.Segment

	// This is a simplified parser - you'd need proper XML parsing
	re := regexp.MustCompile(`<text start="([^"]+)" dur="([^"]+)">([^<]+)</text>`)
	matches := re.FindAllSubmatch(data, -1)

	for _, match := range matches {
		if len(match) == 4 {
			var start, dur float64
			fmt.Sscanf(string(match[1]), "%f", &start)
			fmt.Sscanf(string(match[2]), "%f", &dur)

			segments = append(segments, models.Segment{
				Start: start,
				End:   start + dur,
				Text:  string(match[3]),
			})
		}
	}

	return segments, nil
}

func (y *YouTubePlatform) GetDownloadConfig() *DownloadConfig {
	return &DownloadConfig{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Referer:   "https://www.youtube.com/",
		Headers: map[string]string{
			"Accept-Language": "en-US,en;q=0.9",
		},
		ExtraArgs: []string{
			"--no-check-certificates",
		},
	}
}
