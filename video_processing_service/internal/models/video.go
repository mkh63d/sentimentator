package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	VideoID      string         `gorm:"uniqueIndex:idx_video_platform;not null" json:"video_id"`
	Platform     string         `gorm:"uniqueIndex:idx_video_platform;size:45;not null" json:"platform"`
	OriginalURL  string         `gorm:"type:text;not null" json:"original_url"`
	Title        string         `gorm:"size:255" json:"title"`
	Duration     int            `json:"duration"` // in seconds
	AudioURL     string         `gorm:"type:text" json:"audio_url,omitempty"`
	ThumbnailURL string         `gorm:"type:text" json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time      `gorm:"index:idx_videos_created_at,sort:desc" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Transcripts []Transcript `gorm:"foreignKey:VideoID;references:ID" json:"transcripts,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoMetadata struct {
	VideoID      string
	Platform     string
	Title        string
	Duration     int // in seconds
	HasCaptions  bool
	AudioURL     string
	ThumbnailURL string
}
