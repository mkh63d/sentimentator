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
	"path/filepath"

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
		apiKey:    whisperAPIKey,
		client:    &http.Client{},
		extractor: extractorService,
	}
}

func (w *WhisperStrategy) Name() string {
	return "whisper"
}

func (w *WhisperStrategy) GetTranscript(ctx context.Context, videoURL, platformName, videoID, lang string) (*Result, error) {
	audioResult, err := w.extractor.ExtractAudio(ctx, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract audio from %s: %w", platformName, err)
	}
	audioPath := audioResult.FilePath
	defer os.Remove(audioPath)

	transcript, err := w.transcribeWithWhisper(ctx, audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	return &Result{
		Text:     transcript.Text,
		Segments: transcript.Segments,
		Method:   transcript.Method,
	}, nil
}

// transcribeWithWhisper sends the audio file to OpenAI's Whisper API for transcription.
// This method is platform-agnostic and only cares about the audio file path.
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
		Method:   models.MethodWhisper,
	}, nil
}
