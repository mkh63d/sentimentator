package extractor

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

func checkYtDlpInstalled() bool {
	cmd := exec.Command("yt-dlp", "--version")
	return cmd.Run() == nil
}

// TestYouTubePlatform_ExtractVideoID tests YouTube video ID extraction from various URL formats
func TestYouTubePlatform_ExtractVideoID(t *testing.T) {
	platform := NewYouTubePlatform(t.TempDir())

	tests := []struct {
		name    string
		url     string
		wantID  string
		wantErr bool
	}{
		{
			name:    "Standard watch URL",
			url:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Standard watch URL without www",
			url:     "https://youtube.com/watch?v=dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Short URL (youtu.be)",
			url:     "https://youtu.be/dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Embed URL",
			url:     "https://www.youtube.com/embed/dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "V format URL",
			url:     "https://www.youtube.com/v/dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "URL with additional parameters",
			url:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=10s",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "URL with playlist parameter",
			url:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantID:  "",
			wantErr: true,
		},
		{
			name:    "Invalid YouTube URL",
			url:     "https://www.youtube.com/invalid",
			wantID:  "",
			wantErr: true,
		},
		{
			name:    "Non-YouTube URL",
			url:     "https://www.vimeo.com/123456",
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

// TestYouTubePlatform_SupportsURL tests YouTube URL detection
func TestYouTubePlatform_SupportsURL(t *testing.T) {
	platform := NewYouTubePlatform(t.TempDir())

	tests := []struct {
		name        string
		url         string
		wantSupport bool
	}{
		{
			name:        "YouTube standard URL",
			url:         "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantSupport: true,
		},
		{
			name:        "YouTube short URL",
			url:         "https://youtu.be/dQw4w9WgXcQ",
			wantSupport: true,
		},
		{
			name:        "YouTube embed URL",
			url:         "https://www.youtube.com/embed/dQw4w9WgXcQ",
			wantSupport: true,
		},
		{
			name:        "Vimeo URL",
			url:         "https://vimeo.com/123456",
			wantSupport: false,
		},
		{
			name:        "Instagram URL",
			url:         "https://www.instagram.com/p/ABC123/",
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

// TestYouTubePlatform_GetMetadata tests metadata extraction from real YouTube videos
func TestYouTubePlatform_GetMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	platform := NewYouTubePlatform(t.TempDir())

	tests := []struct {
		name        string
		videoID     string
		wantErr     bool
		expectEmpty bool
	}{
		{
			name:        "Valid video ID",
			videoID:     "jNQXAC9IVRw",
			wantErr:     false,
			expectEmpty: false,
		},
		{
			name:        "Invalid video ID returns empty metadata",
			videoID:     "invalidVideoID123",
			wantErr:     false,
			expectEmpty: true,
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

			// For invalid video IDs, YouTube returns empty metadata rather than error
			if tt.expectEmpty {
				if metadata.Title != "" {
					t.Errorf("expected empty title for invalid video, got: %s", metadata.Title)
				}
				t.Logf("Invalid video correctly returned empty metadata")
				return
			}

			// Verify metadata fields for valid videos
			if metadata.VideoID != tt.videoID {
				t.Errorf("VideoID = %v, want %v", metadata.VideoID, tt.videoID)
			}

			if metadata.Platform != "youtube_platform" {
				t.Errorf("Platform = %v, want youtube_platform", metadata.Platform)
			}

			if metadata.Title == "" {
				t.Error("Title is empty")
			}

			if metadata.OriginalURL == "" {
				t.Error("OriginalURL is empty")
			}

			t.Logf("Metadata extracted: Title=%s, HasCaptions=%v", metadata.Title, metadata.HasCaptions)
		})
	}
}

// TestYouTubePlatform_ExtractAudio tests audio extraction from real YouTube videos
func TestYouTubePlatform_ExtractAudio(t *testing.T) {
	if !checkYtDlpInstalled() {
		t.Skip("yt-dlp is not installed, skipping integration test")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	platform := NewYouTubePlatform(tempDir)

	tests := []struct {
		name    string
		videoID string
		wantErr bool
	}{
		{
			name:    "Extract audio from short video",
			videoID: "jNQXAC9IVRw",
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

// TestYouTubePlatform_Name tests platform name
func TestYouTubePlatform_Name(t *testing.T) {
	platform := NewYouTubePlatform(t.TempDir())

	name := platform.Name()
	expectedName := "youtube_platform"

	if name != expectedName {
		t.Errorf("Name() = %v, want %v", name, expectedName)
	}
}

// BenchmarkYouTubePlatform_ExtractVideoID benchmarks video ID extraction
func BenchmarkYouTubePlatform_ExtractVideoID(b *testing.B) {
	platform := NewYouTubePlatform(b.TempDir())
	testURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := platform.ExtractVideoID(testURL)
		if err != nil {
			b.Fatalf("ExtractVideoID() error = %v", err)
		}
	}
}

// BenchmarkYouTubePlatform_SupportsURL benchmarks URL detection
func BenchmarkYouTubePlatform_SupportsURL(b *testing.B) {
	platform := NewYouTubePlatform(b.TempDir())
	testURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = platform.SupportsURL(testURL)
	}
}
