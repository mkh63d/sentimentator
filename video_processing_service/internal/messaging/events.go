package messaging

import "time"

const (
	EventTypeVideoReceived    = "video.received"
	EventTypeVideoProcessed   = "video.processed"
	EventTypeTranscriptReady  = "transcript.ready"
	EventTypeProcessingFailed = "processing.failed"
)

const (
	QueueVideoProcessing = "video.processing"
	QueueTranscript      = "video.transcript"
	QueueDeadLetter      = "video.dlq"
)

const (
	ExchangeVideo = "video"
	ExchangeTopic = "sentimentator.topic"
)

const (
	RoutingKeyVideoReceived    = "video.received"
	RoutingKeyVideoProcessed   = "video.processed"
	RoutingKeyTranscriptReady  = "transcript.ready"
	RoutingKeyProcessingFailed = "processing.failed"
)

// VideoReceivedEvent represents a new video that needs processing
type VideoReceivedEvent struct {
	VideoID   string    `json:"video_id"`
	URL       string    `json:"url"`
	Platform  string    `json:"platform"`
	UserID    string    `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// VideoProcessedEvent represents a successfully processed video
type VideoProcessedEvent struct {
	VideoID     string    `json:"video_id"`
	URL         string    `json:"url"`
	AudioPath   string    `json:"audio_path,omitempty"`
	Duration    float64   `json:"duration"`
	ProcessedAt time.Time `json:"processed_at"`
}

// TranscriptReadyEvent represents a completed transcription
type TranscriptReadyEvent struct {
	VideoID      string    `json:"video_id"`
	TranscriptID string    `json:"transcript_id"`
	Text         string    `json:"text"`
	Language     string    `json:"language,omitempty"`
	Duration     float64   `json:"duration"`
	WordCount    int       `json:"word_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// ProcessingFailedEvent represents a processing failure
type ProcessingFailedEvent struct {
	VideoID  string    `json:"video_id"`
	Stage    string    `json:"stage"` // e.g., "download", "transcription"
	Error    string    `json:"error"`
	Attempts int       `json:"attempts"`
	FailedAt time.Time `json:"failed_at"`
}
