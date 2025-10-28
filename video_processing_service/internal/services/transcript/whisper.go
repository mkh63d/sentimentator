package transcript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type WhisperStrategy struct {
	apiKey  string
	client  *http.Client
	tempDir string
}

func NewWhisperStrategy(apiKey string) *WhisperStrategy {
	return &WhisperStrategy{
		apiKey:  apiKey,
		client:  &http.Client{},
		tempDir: os.TempDir(),
	}
}

func (w *WhisperStrategy) Name() string {
	return "whisper"
}

func (w *WhisperStrategy) GetTranscript(ctx context.Context, videoURL, platformName, videoID, lang string) (*Result, error) {
	audioPath, err := w.downloadAudio(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download audio: %w", err)
	}
	defer os.Remove(audioPath)

	transcript, err := w.transcribeWithWhisper(ctx, audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	return &Result{
		Text:     transcript.Text,
		Segments: nil,
		Method:   models.MethodWhisper,
	}, nil
}

func (w *WhisperStrategy) downloadAudio(ctx context.Context, videoURL string) (string, error) {
	outputPath := filepath.Join(w.tempDir, fmt.Sprintf("temp_audio_%d.mp3", os.Getpid()))

	cmd := exec.CommandContext(
		ctx, "yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"-o", outputPath,
		"--no-check-certificates",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"--referer", "https://www.youtube.com/",
		"--add-header", "Accept-Language:en-US,en;q=0.9",
		videoURL,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp command failed: %w, stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("audio file was not created at %s", outputPath)
	}

	return outputPath, nil
}

func (w *WhisperStrategy) transcribeWithWhisper(ctx context.Context, audioPath string) (*Result, error) {
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	writer.WriteField("model", "whisper-1")
	writer.WriteField("response_format", "verbose_json")

	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("whisper API error: %s", string(bodyBytes))
	}

	var response struct {
		Text     string `json:"text"`
		Segments []struct {
			Start float64 `json:"start"`
			End   float64 `json:"end"`
			Text  string  `json:"text"`
		} `json:"segments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	segments := make([]models.Segment, len(response.Segments))
	for i, seg := range response.Segments {
		segments[i] = models.Segment{
			Start: seg.Start,
			End:   seg.End,
			Text:  seg.Text,
		}
	}

	return &Result{
		Text:     response.Text,
		Segments: segments,
	}, nil
}
