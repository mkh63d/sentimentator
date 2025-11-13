package messaging

import (
	"context"
	"fmt"
	"log/slog"
)

type Producer struct {
	rabbitmq *RabbitMQ
}

func NewProducer(rabbitmq *RabbitMQ) *Producer {
	return &Producer{
		rabbitmq: rabbitmq,
	}
}

func (p *Producer) PublishVideoProcessed(ctx context.Context, event VideoTranscriptStartedEvent) error {
	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoTranscriptStarted, event)
	if err != nil {

		return fmt.Errorf("failed to publish video transcript started event: %w", err)

	}

	slog.Info("Published video transcript started event",
		slog.String("video_id", event.VideoID),
		slog.String("stage", event.Stage))

	return nil

}

func (p *Producer) PublishTranscriptReady(ctx context.Context, event VideoTranscriptCompletedEvent) error {

	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoTranscriptCompleted, event)
	if err != nil {

		return fmt.Errorf("failed to publish transcript completed event: %w", err)

	}

	slog.Info("Published transcript completed event",
		slog.String("video_id", event.VideoID),

		slog.String("transcript_id", event.TranscriptID))

	return nil

}

func (p *Producer) PublishProcessingFailed(ctx context.Context, event VideoTranscriptFailedEvent) error {

	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoTranscriptFailed, event)
	if err != nil {

		return fmt.Errorf("failed to publish transcript failed event: %w", err)

	}

	slog.Error("Published transcript failed event",
		slog.String("video_id", event.VideoID),

		slog.String("stage", event.Stage),

		slog.String("error", event.Error))

	return nil

}

func (p *Producer) PublishVideoReceived(ctx context.Context, event VideoTranscriptRequestedEvent) error {

	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoTranscriptRequested, event)
	if err != nil {

		return fmt.Errorf("failed to publish transcript requested event: %w", err)
	}

	slog.Info("Published transcript requested event",
		slog.String("analysis_id", event.AnalysisID),
		slog.String("url", event.URL))

	return nil

}
