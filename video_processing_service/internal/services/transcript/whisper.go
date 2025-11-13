package transcript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
)

type WhisperStrategy struct {
	apiKey    string
	client    *http.Client
	extractor extractor.Service
}

func NewWhisperStrategy(extractorService extractor.Service) *WhisperStrategy {
	whisperAPIKey := os.Getenv("OPENAI_API_KEY")
	if whisperAPIKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	return &WhisperStrategy{
		apiKey: whisperAPIKey,
		client: &http.Client{
			Timeout: 15 * time.Minute,
		},
		extractor: extractorService,
	}
}

func (w *WhisperStrategy) Name() string {
	return "whisper"
}

func (w *WhisperStrategy) GetTranscript(ctx context.Context, videoURL, platformName, videoID, lang string) (*Result, error) {
	totalStart := time.Now()

	slog.Info("Starting audio extraction",
		slog.String("platform", platformName),
		slog.String("video_id", videoID))

	extractStart := time.Now()
	audioResult, err := w.extractor.ExtractAudio(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract audio from %s: %w", platformName, err)
	}
	audioPath := audioResult.FilePath
	defer os.Remove(audioPath)

	audioExtractionDuration := time.Since(extractStart).Seconds()

	slog.Info("Audio extracted successfully",
		slog.String("video_id", videoID),
		slog.String("file_path", audioPath),
		slog.Int64("file_size_bytes", audioResult.FileSize),
		slog.Int64("duration_seconds", int64(audioExtractionDuration)))

	slog.Info("Starting Whisper transcription",
		slog.String("video_id", videoID),
		slog.Int64("file_size_mb", audioResult.FileSize/(1024*1024)))

	transcribeStart := time.Now()
	transcript, err := w.transcribeWithWhisper(ctx, audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	transcriptionDuration := time.Since(transcribeStart).Seconds()
	totalDuration := time.Since(totalStart).Seconds()

	slog.Info("Whisper transcription completed",
		slog.String("video_id", videoID),
		slog.Int("text_length", len(transcript.Text)),
		slog.Int("segments_count", len(transcript.Segments)),
		slog.Int64("transcription_duration_seconds", int64(transcriptionDuration)),
		slog.Int64("total_duration_seconds", int64(totalDuration)))

	return &Result{
		Text:                    transcript.Text,
		Segments:                transcript.Segments,
		Method:                  transcript.Method,
		DetectedLanguage:        transcript.DetectedLanguage,
		AudioExtractionDuration: int64(audioExtractionDuration),
		TranscriptionDuration:   int64(transcriptionDuration),
		TotalProcessingDuration: int64(totalDuration),
		AudioFileSizeBytes:      audioResult.FileSize,
	}, nil
}

// transcribeWithWhisper sends the audio file to OpenAI's Whisper API for transcription.
// This method is platform-agnostic and only cares about the audio file path.
func (w *WhisperStrategy) transcribeWithWhisper(ctx context.Context, audioPath string) (*Result, error) {
	fileInfo, err := os.Stat(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat audio file: %w", err)
	}

	slog.Info("Preparing audio file for Whisper API",
		slog.String("file_path", audioPath),
		slog.Int64("file_size_mb", fileInfo.Size()/(1024*1024)))

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

	slog.Info("Sending request to Whisper API",
		slog.Int("request_size_mb", body.Len()/(1024*1024)))

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	apiCallStart := time.Now()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("Received response from Whisper API",
		slog.Int("status_code", resp.StatusCode),
		slog.String("api_call_duration", time.Since(apiCallStart).String()))

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("Whisper API error response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response_body", string(bodyBytes)))
		return nil, fmt.Errorf("whisper API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Text     string `json:"text"`
		Language string `json:"language"`
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

	slog.Info("Whisper API response decoded successfully",
		slog.Int("text_chars", len(response.Text)),
		slog.Int("segments", len(segments)))

	return &Result{
		Text:             response.Text,
		Segments:         segments,
		Method:           models.MethodWhisper,
		DetectedLanguage: response.Language,
	}, nil
}
