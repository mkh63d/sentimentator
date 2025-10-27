package models

type Video struct {
	ID        string `json:"id" db:"id"`
	URL       string `json:"url" db:"url"`
	Platform  string `json:"platform" db:"platform"`
	VideoID   string `json:"video_id" db:"video_id"`
	Title     string `json:"title" db:"title"`
	Duration  int    `json:"duration" db:"duration"` //  in seconds
	AudioURL  string `json:"audio_url,omitempty" db:"audio_url"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type VideoMetadata struct {
	VideoID     string
	Platform    string
	Title       string
	Duration    int // in seconds
	HasCaptions bool
	AudioURL    string
}
