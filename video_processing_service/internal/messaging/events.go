package messaging

import "time"

const (
	EventTypeVideoTranscriptRequested = "video.transcript.requested"
	EventTypeVideoTranscriptStarted   = "video.transcript.started"
	EventTypeVideoTranscriptCompleted = "video.transcript.completed"
	EventTypeVideoTranscriptFailed    = "video.transcript.failed"
)

const (
	QueueVideoProcessing = "video.processing"
	QueueTranscript      = "video.transcript"
)

const (
	ExchangeVideo = "video"
	ExchangeTopic = "sentimentator.topic"
)

const (
	RoutingKeyVideoTranscriptRequested = "video.transcript.requested"
	RoutingKeyVideoTranscriptStarted   = "video.transcript.started"
	RoutingKeyVideoTranscriptCompleted = "video.transcript.completed"
	RoutingKeyVideoTranscriptFailed    = "video.transcript.failed"
)

// VideoTranscriptRequestedEvent represents a new video URL that needs processing to get a transcript
type VideoTranscriptRequestedEvent struct {
	AnalysisID string    `json:"analysis_id"`
	URL        string    `json:"url"`
	Timestamp  time.Time `json:"timestamp"`
}

// VideoTranscriptStartedEvent represents that processing for the transcript has begun
type VideoTranscriptStartedEvent struct {
	AnalysisID string    `json:"analysis_id"`
	VideoID    string    `json:"video_id,omitempty"`
	Stage      string    `json:"stage,omitempty"` // e.g., "download", "transcription"
	StartedAt  time.Time `json:"started_at"`
}

// VideoTranscriptCompletedEvent represents a finished transcription
type VideoTranscriptCompletedEvent struct {
	AnalysisID   string    `json:"analysis_id"`
	VideoID      string    `json:"video_id,omitempty"`
	TranscriptID string    `json:"transcript_id,omitempty"`
	Text         string    `json:"text"`
	Language     string    `json:"language,omitempty"`
	Duration     float64   `json:"duration,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// VideoTranscriptFailedEvent represents a failure in the transcript workflow.
type VideoTranscriptFailedEvent struct {
	AnalysisID string    `json:"analysis_id"`
	VideoID    string    `json:"video_id,omitempty"`
	Stage      string    `json:"stage"` // "download" | "transcription" | "storage" ...
	Error      string    `json:"error"`
	Attempts   int       `json:"attempts"`
	FailedAt   time.Time `json:"failed_at"`
}
