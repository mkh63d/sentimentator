package utils

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

type VideoInfo struct {
	VideoID     string
	Platform    string
	OriginalURL string
}

// ParseVideoURL parses a video URL and extracts the video ID and platform
func ParseVideoURL(rawURL string) (*VideoInfo, error) {
	if rawURL == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}

	host := strings.ToLower(parsedURL.Host)

	host = strings.TrimPrefix(host, "www.")

	// YouTube URL patterns
	if strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be") {
		videoID, err := extractYouTubeID(parsedURL, rawURL)
		if err != nil {
			return nil, err
		}
		return &VideoInfo{
			VideoID:     videoID,
			Platform:    "youtube",
			OriginalURL: rawURL,
		}, nil
	}

	return nil, errors.New("unsupported video platform")
}

// extractYouTubeID extracts the video ID from a YouTube URL
func extractYouTubeID(parsedURL *url.URL, rawURL string) (string, error) {
	host := strings.ToLower(parsedURL.Host)
	host = strings.TrimPrefix(host, "www.")

	// youtu.be short links
	if host == "youtu.be" {
		videoID := strings.TrimPrefix(parsedURL.Path, "/")
		if videoID != "" {
			return videoID, nil
		}
	}

	// youtube.com URLs
	if strings.Contains(host, "youtube.com") {
		// Check query parameter 'v'
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

	return "", errors.New("could not extract YouTube video ID")
}
