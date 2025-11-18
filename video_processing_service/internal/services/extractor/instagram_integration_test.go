package extractor

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

func checkYtDlpInstalledForInstagram() bool {
	cmd := exec.Command("yt-dlp", "--version")
	return cmd.Run() == nil
}

// TestInstagramPlatform_ExtractVideoID tests Instagram video ID extraction from various URL formats
func TestInstagramPlatform_ExtractVideoID(t *testing.T) {
	platform := NewInstagramPlatform(t.TempDir())

	tests := []struct {
		name    string
		url     string
		wantID  string
		wantErr bool
	}{
		{
			name:    "Instagram reels URL",
			url:     "https://www.instagram.com/reels/DQg6qZyjLTp/",
			wantID:  "DQg6qZyjLTp",
			wantErr: false,
		},
		{
			name:    "Instagram reels URL without trailing slash",
			url:     "https://www.instagram.com/reels/DQg6qZyjLTp",
			wantID:  "DQg6qZyjLTp",
			wantErr: false,
		},
		{
			name:    "Instagram post URL (p format)",
			url:     "https://www.instagram.com/p/ABC123xyz/",
			wantID:  "ABC123xyz",
			wantErr: false,
		},
		{
			name:    "Instagram TV URL (tv format)",
			url:     "https://www.instagram.com/tv/DEF456uvw/",
			wantID:  "DEF456uvw",
			wantErr: false,
		},
		{
			name:    "Instagram URL without www",
			url:     "https://instagram.com/reels/GHI789rst/",
			wantID:  "GHI789rst",
			wantErr: false,
		},
		{
			name:    "Instagram short URL (instagr.am)",
			url:     "https://instagr.am/reels/JKL012opq/",
			wantID:  "JKL012opq",
			wantErr: false,
		},
		{
			name:    "Instagram URL with query parameters",
			url:     "https://www.instagram.com/reels/MNO345lmn/?utm_source=ig_web_copy_link",
			wantID:  "MNO345lmn",
			wantErr: false,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantID:  "",
			wantErr: true,
		},
		{
			name:    "Invalid Instagram URL",
			url:     "https://www.instagram.com/invalid",
			wantID:  "",
			wantErr: true,
		},
		{
			name:    "Non-Instagram URL",
			url:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoID, err := platform.ExtractVideoID(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("ExtractVideoID() error = %v", err)
			}

			if videoID != tt.wantID {
				t.Errorf("ExtractVideoID() = %v, want %v", videoID, tt.wantID)
			}
		})
	}
}

// TestInstagramPlatform_SupportsURL tests Instagram URL detection
func TestInstagramPlatform_SupportsURL(t *testing.T) {
	platform := NewInstagramPlatform(t.TempDir())

	tests := []struct {
		name        string
		url         string
		wantSupport bool
	}{
		{
			name:        "Instagram reels URL",
			url:         "https://www.instagram.com/reels/DQg6qZyjLTp/",
			wantSupport: true,
		},
		{
			name:        "Instagram post URL",
			url:         "https://www.instagram.com/p/ABC123/",
			wantSupport: true,
		},
		{
			name:        "Instagram TV URL",
			url:         "https://www.instagram.com/tv/DEF456/",
			wantSupport: true,
		},
		{
			name:        "Instagram short URL",
			url:         "https://instagr.am/reels/GHI789/",
			wantSupport: true,
		},
		{
			name:        "YouTube URL",
			url:         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantSupport: false,
		},
		{
			name:        "Vimeo URL",
			url:         "https://vimeo.com/123456",
			wantSupport: false,
		},
		{
			name:        "Empty URL",
			url:         "",
			wantSupport: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supports := platform.SupportsURL(tt.url)

			if supports != tt.wantSupport {
				t.Errorf("SupportsURL() = %v, want %v", supports, tt.wantSupport)
			}
		})
	}
}

// TestInstagramPlatform_GetMetadata tests metadata extraction from real Instagram videos
func TestInstagramPlatform_GetMetadata(t *testing.T) {
	if !checkYtDlpInstalledForInstagram() {
		t.Skip("yt-dlp is not installed, skipping integration test")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	platform := NewInstagramPlatform(t.TempDir())

	tests := []struct {
		name        string
		videoID     string
		wantErr     bool
		expectEmpty bool
	}{
		{
			name:        "Valid Instagram reel",
			videoID:     "DQg6qZyjLTp",
			wantErr:     false,
			expectEmpty: false,
		},
		{
			name:        "Invalid video ID",
			videoID:     "invalidVideoID123",
			wantErr:     true,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			metadata, err := platform.GetMetadata(ctx, tt.videoID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetMetadata() error = %v", err)
			}

			// Verify metadata fields
			if metadata.VideoID != tt.videoID {
				t.Errorf("VideoID = %v, want %v", metadata.VideoID, tt.videoID)
			}

			if metadata.Platform != "instagram_platform" {
				t.Errorf("Platform = %v, want instagram_platform", metadata.Platform)
			}

			if metadata.Title == "" {
				t.Logf("Warning: Title is empty for video ID %s", tt.videoID)
			}

			if metadata.OriginalURL == "" {
				t.Logf("Warning: OriginalURL is empty for video ID %s", tt.videoID)
			}

			t.Logf("Metadata extracted: Title=%s, Duration=%d, HasCaptions=%v",
				metadata.Title, metadata.Duration, metadata.HasCaptions)
		})
	}
}

// TestInstagramPlatform_ExtractAudio tests audio extraction from real Instagram videos
func TestInstagramPlatform_ExtractAudio(t *testing.T) {
	if !checkYtDlpInstalledForInstagram() {
		t.Skip("yt-dlp is not installed, skipping integration test")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	platform := NewInstagramPlatform(tempDir)

	tests := []struct {
		name    string
		videoID string
		wantErr bool
	}{
		{
			name:    "Extract audio from Instagram reel",
			videoID: "DQg6qZyjLTp",
			wantErr: false,
		},
		{
			name:    "Invalid video ID",
			videoID: "invalidVideoID123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			result, err := platform.ExtractAudio(ctx, tt.videoID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("ExtractAudio() error = %v", err)
			}

			defer func() {
				if result.FilePath != "" {
					os.Remove(result.FilePath)
				}
			}()

			// Verify the audio file was created
			if result.FilePath == "" {
				t.Error("FilePath is empty")
			}

			fileInfo, err := os.Stat(result.FilePath)
			if err != nil {
				t.Fatalf("audio file does not exist: %v", err)
			}

			// Verify file has content (should be at least 1KB for even short videos)
			if fileInfo.Size() < 1024 {
				t.Errorf("audio file is too small: %d bytes", fileInfo.Size())
			}

			// Verify format
			if result.Format != "mp3" {
				t.Errorf("expected format mp3, got %s", result.Format)
			}

			t.Logf("Successfully extracted audio: %s (size: %d bytes)", result.FilePath, fileInfo.Size())
		})
	}
}

// TestInstagramPlatform_Name tests platform name
func TestInstagramPlatform_Name(t *testing.T) {
	platform := NewInstagramPlatform(t.TempDir())

	name := platform.Name()
	expectedName := "instagram_platform"

	if name != expectedName {
		t.Errorf("Name() = %v, want %v", name, expectedName)
	}
}

// BenchmarkInstagramPlatform_ExtractVideoID benchmarks video ID extraction
func BenchmarkInstagramPlatform_ExtractVideoID(b *testing.B) {
	platform := NewInstagramPlatform(b.TempDir())
	testURL := "https://www.instagram.com/reels/DQg6qZyjLTp/"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := platform.ExtractVideoID(testURL)
		if err != nil {
			b.Fatalf("ExtractVideoID() error = %v", err)
		}
	}
}

// BenchmarkInstagramPlatform_SupportsURL benchmarks URL detection
func BenchmarkInstagramPlatform_SupportsURL(b *testing.B) {
	platform := NewInstagramPlatform(b.TempDir())
	testURL := "https://www.instagram.com/reels/DQg6qZyjLTp/"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = platform.SupportsURL(testURL)
	}
}

// TestInstagramPlatform_Integration tests full integration with real Instagram URLs
func TestInstagramPlatform_Integration(t *testing.T) {
	if !checkYtDlpInstalledForInstagram() {
		t.Skip("yt-dlp is not installed, skipping integration test")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	platform := NewInstagramPlatform(tempDir)

	// Test with multiple Instagram URLs from VIDEO_LINKS.md
	urls := []string{
		"https://www.instagram.com/reels/DQg6qZyjLTp/",
		"https://www.instagram.com/reels/DN6h1U8D5Oe/",
		"https://www.instagram.com/reels/DQQhHlUjLBm/",
	}

	for _, testURL := range urls {
		t.Run("Integration_"+testURL, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			videoID, err := platform.ExtractVideoID(testURL)
			if err != nil {
				t.Fatalf("ExtractVideoID() error = %v", err)
			}

			// Get metadata
			metadata, err := platform.GetMetadata(ctx, videoID)
			if err != nil {
				t.Logf("Warning: GetMetadata failed for %s: %v", testURL, err)
				// Don't fail on metadata extraction as it might be restricted
				return
			}

			if metadata.VideoID != videoID {
				t.Errorf("Metadata VideoID mismatch: got %v, want %v", metadata.VideoID, videoID)
			}

			if metadata.Platform != "instagram_platform" {
				t.Errorf("Platform mismatch: got %v, want instagram_platform", metadata.Platform)
			}

			t.Logf("Integration test passed for %s: Title=%s, Duration=%d", testURL, metadata.Title, metadata.Duration)
		})
	}
}
