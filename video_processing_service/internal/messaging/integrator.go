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

// PublishVideoFromModel publishes a transcript request event from a Video model

func (si *ServiceIntegrator) PublishVideoFromModel(ctx context.Context, video *models.Video) error {
	event := VideoTranscriptRequestedEvent{

		AnalysisID: fmt.Sprintf("%d", video.ID),

		URL: video.OriginalURL,

		Timestamp: time.Now(),
	}

	return si.producer.PublishVideoReceived(ctx, event)

}

// PublishTranscriptFromModel publishes a transcript completed event from Transcript and Video models
func (si *ServiceIntegrator) PublishTranscriptFromModel(ctx context.Context, transcript *models.Transcript, video *models.Video) error {
	event := VideoTranscriptCompletedEvent{

		AnalysisID: fmt.Sprintf("%d", video.ID),

		VideoID: fmt.Sprintf("%d", video.ID),

		TranscriptID: fmt.Sprintf("%d", transcript.ID),

		Text: transcript.Text,

		Language: "en",

		Duration: video.Duration,

		CreatedAt: transcript.CreatedAt,
	}

	return si.producer.PublishTranscriptReady(ctx, event)

}
