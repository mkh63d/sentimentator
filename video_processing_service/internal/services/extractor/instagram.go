package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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

	cmd := exec.CommandContext(ctx, "yt-dlp", "-J", "--no-warning", "--skip-download", videoURL)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(out, &obj); err != nil {
		return nil, err
	}

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

	metadata := &VideoMetadata{
		VideoID:                   videoID,
		Platform:                  ip.Name(),
		Title:                     title,
		Duration:                  duration,
		HasCaptions:               hasCaptions,
		AvailableCaptionLanguages: availableLangs,
		OriginalURL:               origURL,
	}

	return metadata, nil
}

func (ip *InstagramPlatform) GetCaptions(ctx context.Context, videoID, lang string) ([]Segment, error) {
	return make([]Segment, 0), nil
}

func (ip *InstagramPlatform) ExtractAudio(ctx context.Context, videoID string) (*AudioResult, error) {
	outputPath := filepath.Join(ip.tempDir, fmt.Sprintf("temp_instagram_audio_%d.mp3", os.Getpid()))
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
