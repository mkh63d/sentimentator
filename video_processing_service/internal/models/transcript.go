package models

import "time"

type TranscriptStatus string

const (
	StatusPending    TranscriptStatus = "pending"
	StatusProcessing TranscriptStatus = "processing"
	StatusCompleted  TranscriptStatus = "completed"
	StatusFailed     TranscriptStatus = "failed"
)

type TranscriptMethod string

const (
	MethodCaption TranscriptMethod = "caption"
	MethodWhisper TranscriptMethod = "whisper"
)

type Transcript struct {
	ID          string           `json:"id" db:"id"`
	VideoID     string           `json:"video_id" db:"video_id"`
	Language    string           `json:"language" db:"language"`
	Method      TranscriptMethod `json:"method" db:"method"`
	Status      TranscriptStatus `json:"status" db:"status"`
	Text        string           `json:"text,omitempty" db:"text"`
	Segments    string           `json:"segments,omitempty" db:"-"`
	Error       string           `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty" db:"completed_at"`
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type TranscriptRequest struct {
	VideoURL string `json:"video_url" validate:"required,url"`
	Language string `json:"language,omitempty"`
}

type TranscriptResponse struct {
	ID         string           `json:"id"`
	Status     TranscriptStatus `json:"status"`
	Video      *Video           `json:"video,omitempty"`
	Transcript *Transcript      `json:"transcript,omitempty"`
}
