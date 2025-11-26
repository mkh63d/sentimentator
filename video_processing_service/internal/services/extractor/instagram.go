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

type InstagramPlatform struct {
	client  *http.Client
	tempDir string
}

func NewInstagramPlatform(tempDir string) Platform {
	return &InstagramPlatform{
		client:  &http.Client{},
		tempDir: tempDir}
}

func (ip *InstagramPlatform) Name() string {
	return "instagram_platform"
}

func (ip *InstagramPlatform) SupportsURL(videoURL string) bool {
	return strings.Contains(videoURL, "instagram.com") || strings.Contains(videoURL, "instagr.am")
}

func (ip *InstagramPlatform) ExtractVideoID(rawVideoURL string) (string, error) {
	if rawVideoURL == "" || strings.TrimSpace(rawVideoURL) == "" {
		return "", errors.New("instagram Video URL cannot be empty")
	}

	u, err := url.Parse(rawVideoURL)
	if err != nil {
		return "", err
	}

	//INFO: Expected Instagram URL formats:
	// https://www.instagram.com/reels/jd100%/
	// https://www.instagram.com/p/jd100%/
	// https://www.instagram.com/tv/jd100%/
	re := regexp.MustCompile(`(?:reels|p|tv)\/([^\/]+)`)
	matches := re.FindStringSubmatch(u.Path)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("unable to extract Instagram video ID from URL: %s", rawVideoURL)
}

func (ip *InstagramPlatform) GetMetadata(ctx context.Context, videoID string) (*VideoMetadata, error) {
	if videoID == "" || strings.TrimSpace(videoID) == "" {
		return nil, errors.New("instagram video ID cannot be empty")
	}

	videoURL := fmt.Sprintf("https://www.instagram.com/reels/%s/", videoID)

	args := []string{"-J", "--no-warning", "--skip-download"}

	cookiesPath := os.Getenv("INSTAGRAM_COOKIES_PATH")
	if cookiesPath != "" {
		if _, err := os.Stat(cookiesPath); err == nil {
			args = append(args, "--cookies", cookiesPath)
		}
	}

	args = append(args, videoURL)
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	out, err := cmd.Output()
	if err != nil {
		return ip.getMetadataFromOEmbed(ctx, videoURL, videoID)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(out, &obj); err != nil {
		return nil, err
	}

	// Debug: Log available keys to understand the structure
	slog.DebugContext(ctx, "Instagram metadata keys", slog.Any("keys", getKeys(obj)))

	title := ""
	if t, ok := obj["title"].(string); ok {
		title = t
	}

	duration := 0
	if d, ok := obj["duration"].(float64); ok {
		duration = int(d)
	} else if d, ok := obj["duration"].(int); ok {
		duration = d
	}

	origURL := ""
	if u, ok := obj["webpage_url"].(string); ok {
		origURL = u
	} else if u, ok := obj["url"].(string); ok {
		origURL = u
	}

	hasCaptions := false
	availableLangs := make([]string, 0)

	if subs, ok := obj["subtitles"].(map[string]interface{}); ok && len(subs) > 0 {
		hasCaptions = true
		for lang := range subs {
			availableLangs = append(availableLangs, lang)
		}
	}
	if autos, ok := obj["automatic_captions"].(map[string]interface{}); ok && len(autos) > 0 {
		hasCaptions = true
		for lang := range autos {
			availableLangs = append(availableLangs, lang)
		}
	}

	thumbnailURL := ""
	if thumbnail, ok := obj["thumbnail"].(string); ok && thumbnail != "" {
		thumbnailURL = thumbnail
	}
	if thumbnailURL == "" {
		if thumbnails, ok := obj["thumbnails"].([]interface{}); ok && len(thumbnails) > 0 {
			if lastThumb, ok := thumbnails[len(thumbnails)-1].(map[string]interface{}); ok {
				if url, ok := lastThumb["url"].(string); ok {
					thumbnailURL = url
				}
			}
		}
	}
	if thumbnailURL == "" {
		if thumbObj, ok := obj["thumbnail"].(map[string]interface{}); ok {
			if url, ok := thumbObj["url"].(string); ok {
				thumbnailURL = url
			}
		}
	}

	if thumbnailURL != "" {
		slog.DebugContext(ctx, "Instagram thumbnail extracted", slog.String("url", thumbnailURL))
	} else {
		slog.WarnContext(ctx, "No Instagram thumbnail found in metadata")
	}

	metadata := &VideoMetadata{
		VideoID:                   videoID,
		Platform:                  ip.Name(),
		Title:                     title,
		Duration:                  duration,
		HasCaptions:               hasCaptions,
		AvailableCaptionLanguages: availableLangs,
		OriginalURL:               origURL,
		ThumbnailURL:              thumbnailURL,
	}

	return metadata, nil
}

// getMetadataFromOEmbed fetches metadata from Instagram's oEmbed API as fallback
func (ip *InstagramPlatform) getMetadataFromOEmbed(ctx context.Context, videoURL, videoID string) (*VideoMetadata, error) {
	oembedURL := fmt.Sprintf("https://api.instagram.com/oembed/?url=%s", url.QueryEscape(videoURL))
	req, err := http.NewRequestWithContext(ctx, "GET", oembedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create oEmbed request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	resp, err := ip.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch oEmbed data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oEmbed API returned status %d", resp.StatusCode)
	}

	//INFO: Docs https://developers.facebook.com/docs/instagram-platform/oembed/
	var oembed struct {
		Title        string `json:"title"`
		AuthorName   string `json:"author_name"`
		ThumbnailURL string `json:"thumbnail_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&oembed); err != nil {
		return nil, fmt.Errorf("failed to decode oEmbed response: %w", err)
	}

	title := oembed.Title
	if title == "" && oembed.AuthorName != "" {
		title = fmt.Sprintf("Post by %s", oembed.AuthorName)
	}

	return &VideoMetadata{
		VideoID:                   videoID,
		Platform:                  ip.Name(),
		Title:                     title,
		Duration:                  -1,
		HasCaptions:               false,
		AvailableCaptionLanguages: make([]string, 0),
		OriginalURL:               videoURL,
		ThumbnailURL:              oembed.ThumbnailURL,
	}, nil
}

func (ip *InstagramPlatform) GetCaptions(ctx context.Context, videoID, lang string) ([]Segment, error) {
	return make([]Segment, 0), nil
}

func (ip *InstagramPlatform) ExtractAudio(ctx context.Context, videoID string) (*AudioResult, error) {
	outputPath := filepath.Join(ip.tempDir, fmt.Sprintf("temp_instagram_audio_%s_%d.mp3", videoID, time.Now().UnixNano()))
	videoURL := fmt.Sprintf("https://www.instagram.com/reels/%s", videoID)

	args := []string{
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		"--no-check-certificates",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"--referer", "https://www.instagram.com/",
		"--add-header", "Accept-Language: en-US,en;q=0.9",
		"--geo-bypass",
		"--retries", "2",
		"--fragment-retries", "3",
		"--no-cache-dir",
	}

	cookiesPath := os.Getenv("INSTAGRAM_COOKIES_PATH")
	if cookiesPath != "" {
		if _, err := os.Stat(cookiesPath); err == nil {
			args = append(args, "--cookies", cookiesPath)
		}
	}

	args = append(args, videoURL)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed for %s: %w, stderr: %s",
			ip.Name(), err, stderr.String())
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

// getKeys returns all keys from a map for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
