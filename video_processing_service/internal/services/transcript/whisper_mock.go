package transcript

import (
	"context"
	"math/rand"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

// MockWhisperStrategy implements a mock version of WhisperStrategy for testing
// without making actual OpenAI API calls. It returns Lorem Ipsum as transcript data.
//
// Usage:
// Set WHISPER_MOCK_MODE=true in environment to enable mock mode.
type MockWhisperStrategy struct{}

func NewMockWhisperStrategy() *MockWhisperStrategy {
	return &MockWhisperStrategy{}
}

func (m *MockWhisperStrategy) Name() string {
	return "mock_whisper_api"
}

// GetTranscript returns mock Lorem Ipsum transcript data without making API calls
func (m *MockWhisperStrategy) GetTranscript(ctx context.Context, videoURL, platformName, videoID, lang string) (*Result, error) {
	minSeconds := int64(15)
	maxSeconds := int64(120)

	randomSeconds := minSeconds + rand.Int63n(maxSeconds-minSeconds+1)
	randomDuration := time.Duration(randomSeconds) * time.Second

	select {
	case <-time.After(randomDuration):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	segments := generateMockSegments()

	// Build full text from segments
	var fullText string
	for i, seg := range segments {
		if i > 0 {
			fullText += " "
		}
		fullText += seg.Text
	}

	return &Result{
		Text:                    fullText,
		Segments:                segments,
		Method:                  models.MockMethodWhisper,
		DetectedLanguage:        "latin",
		AudioExtractionDuration: -1,
		TranscriptionDuration:   randomSeconds,
		TotalProcessingDuration: randomSeconds,
		AudioFileSizeBytes:      -1,
	}, nil
}

func generateMockSegments() []models.Segment {
	loremTexts := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
		"Nisi ut aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse.",
		"Cillum dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident sunt in culpa.",
		"Qui officia deserunt mollit anim id est laborum.",
		"Sed ut perspiciatis unde omnis iste natus error sit voluptatem.",
		"Accusantium doloremque laudantium totam rem aperiam.",
		"Eaque ipsa quae ab illo inventore veritatis et quasi architecto.",
		"Beatae vitae dicta sunt explicabo nemo enim ipsam voluptatem.",
		"Quia voluptas sit aspernatur aut odit aut fugit.",
		"Sed quia consequuntur magni dolores eos qui ratione voluptatem.",
		"Neque porro quisquam est qui dolorem ipsum quia dolor sit amet.",
	}

	segments := make([]models.Segment, 0, len(loremTexts))
	currentTime := 0.0

	r := rand.New(rand.NewSource(42))

	for _, text := range loremTexts {
		duration := 3.0 + r.Float64()*2.0

		segments = append(segments, models.Segment{
			Start: currentTime,
			End:   currentTime + duration,
			Text:  text,
		})

		currentTime += duration
		currentTime += 0.2 + r.Float64()*0.3
	}

	return segments
}

func GenerateMockTranscriptText() string {
	return `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris. Nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse. Cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident sunt in culpa. Qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem. Accusantium doloremque laudantium totam rem aperiam. Eaque ipsa quae ab illo inventore veritatis et quasi architecto. Beatae vitae dicta sunt explicabo nemo enim ipsam voluptatem. Quia voluptas sit aspernatur aut odit aut fugit. Sed quia consequuntur magni dolores eos qui ratione voluptatem. Neque porro quisquam est qui dolorem ipsum quia dolor sit amet.`
}
