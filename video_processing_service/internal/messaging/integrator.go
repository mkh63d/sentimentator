package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

type ServiceIntegrator struct {
	producer *Producer
}

func NewServiceIntegrator(producer *Producer) *ServiceIntegrator {
	return &ServiceIntegrator{
		producer: producer,
	}
}

// PublishVideoFromModel publishes a video event from a Video model
func (si *ServiceIntegrator) PublishVideoFromModel(ctx context.Context, video *models.Video) error {
	event := VideoReceivedEvent{
		VideoID:   fmt.Sprintf("%d", video.ID),
		URL:       video.OriginalURL,
		Platform:  video.Platform,
		Timestamp: time.Now(),
	}

	return si.producer.PublishVideoReceived(ctx, event)
}

// PublishTranscriptFromModel publishes a transcript ready event from Transcript and Video models
func (si *ServiceIntegrator) PublishTranscriptFromModel(ctx context.Context, transcript *models.Transcript, video *models.Video) error {
	event := TranscriptReadyEvent{
		VideoID:      fmt.Sprintf("%d", video.ID),
		TranscriptID: fmt.Sprintf("%d", transcript.ID),
		Text:         transcript.Text,
		Language:     "en",
		Duration:     float64(video.Duration),
		WordCount:    len(transcript.Text) / 5, // is there a better way to estimate word count?
		CreatedAt:    transcript.CreatedAt,
	}

	return si.producer.PublishTranscriptReady(ctx, event)
}
