package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

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
	MethodWhisper TranscriptMethod = "whisper_api"
)

type Transcript struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	VideoID     uint             `gorm:"not null;index" json:"video_id"`
	Language    string           `gorm:"size:10;default:'en'" json:"language"`
	Method      TranscriptMethod `gorm:"type:varchar(50);not null" json:"method"`
	Status      TranscriptStatus `gorm:"type:varchar(50);default:'pending';index" json:"status"`
	Text        string           `gorm:"type:text" json:"text,omitempty"`
	Segments    SegmentArray     `gorm:"type:jsonb" json:"segments,omitempty"`
	Error       string           `gorm:"type:text" json:"error,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`

	Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (Transcript) TableName() string {
	return "transcripts"
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type SegmentArray []Segment

func (s *SegmentArray) Scan(value any) error {
	if value == nil {
		*s = SegmentArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan SegmentArray")
	}

	return json.Unmarshal(bytes, s)
}

func (s SegmentArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}
