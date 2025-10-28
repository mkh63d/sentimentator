package queue

// TranscriptJobPayload represents the payload for a transcript job
type TranscriptJobPayload struct {
	VideoURL     string `json:"video_url"`
	Language     string `json:"language"`
	VideoID      uint   `json:"video_id"`
	TranscriptID uint   `json:"transcript_id"`
	JobID        string `json:"job_id"`
}

// NOT IMPLEMENTED
type AudioExtractJobPayload struct {
	VideoURL string `json:"video_url"`
	VideoID  uint   `json:"video_id"`
	Platform string `json:"platform"`
}
