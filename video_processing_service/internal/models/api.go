package models

// TranscriptRequest represents the request body for transcript creation
type TranscriptRequest struct {
	VideoURL string `json:"video_url" binding:"required"`
	Language string `json:"language,omitempty"`
}

// TranscriptResponse represents the response for transcript requests
type TranscriptResponse struct {
	ID     string           `json:"id"`
	Status TranscriptStatus `json:"status"`
	Error  string           `json:"error,omitempty"`
}

// VideoProcessRequest represents the request for video processing
type VideoProcessRequest struct {
	VideoURL string `json:"video_url" binding:"required"`
}

// VideoProcessResponse represents the response for video processing
type VideoProcessResponse struct {
	VideoID     string `json:"video_id"`
	Platform    string `json:"platform"`
	DatabaseID  uint   `json:"database_id"`
	Status      string `json:"status"`
	Title       string `json:"title,omitempty"`
	OriginalURL string `json:"original_url"`
}
