package video

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type InstagramPlatform struct {
	client *http.Client
}

func NewInstagramPlatform() *InstagramPlatform {
	return &InstagramPlatform{
		client: &http.Client{},
	}
}

func (i *InstagramPlatform) Name() string {
	return "instagram"
}

func (i *InstagramPlatform) SupportsURL(videoURL string) bool {
	return strings.Contains(videoURL, "instagram.com")
}

func (i *InstagramPlatform) ExtractVideoID(videoURL string) (string, error) {
	re := regexp.MustCompile(`/(?:p|reel|tv)/([^/?]+)`)
	matches := re.FindStringSubmatch(videoURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Instagram URL")
	}
	return matches[1], nil
}

func (i *InstagramPlatform) GetMetadata(
	ctx context.Context,
	videoID string,
) (*models.VideoMetadata, error) {
	url := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1&__d=dis", videoID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []struct {
			VideoURL string `json:"video_url"`
			Title    struct {
				Text string `json:"text"`
			} `json:"edge_media_to_caption"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse Instagram data: %w", err)
	}

	return &models.VideoMetadata{
		VideoID:     videoID,
		Platform:    "instagram",
		HasCaptions: false, // Instagram doesn't provide captions
	}, nil
}

func (i *InstagramPlatform) GetCaptions(
	ctx context.Context,
	videoID, lang string,
) ([]models.Segment, error) {
	return nil, fmt.Errorf("Instagram does not provide captions")
}

func (i *InstagramPlatform) GetAudioURL(
	ctx context.Context,
	videoID string,
) (string, error) {
	// TODO: implement audio extraction for Instagram
	return fmt.Sprintf("https://www.instagram.com/p/%s/", videoID), nil
}

func (i *InstagramPlatform) GetDownloadConfig() *DownloadConfig {
	return &DownloadConfig{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Referer:   "https://www.instagram.com/",
		Headers: map[string]string{
			"Accept-Language":  "en-US,en;q=0.9",
			"X-Requested-With": "XMLHttpRequest",
		},
		ExtraArgs: []string{
			"--no-check-certificates",
		},
	}
}
