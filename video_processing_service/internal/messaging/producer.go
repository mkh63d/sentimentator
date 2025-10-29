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

func (p *Producer) PublishVideoProcessed(ctx context.Context, event VideoProcessedEvent) error {
	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoProcessed, event)
	if err != nil {
		return fmt.Errorf("failed to publish video processed event: %w", err)
	}

	slog.Info("Published video processed event", slog.String("video_id", event.VideoID))
	return nil
}

func (p *Producer) PublishTranscriptReady(ctx context.Context, event TranscriptReadyEvent) error {
	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyTranscriptReady, event)
	if err != nil {
		return fmt.Errorf("failed to publish transcript ready event: %w", err)
	}

	slog.Info("Published transcript ready event",
		slog.String("video_id", event.VideoID),
		slog.String("transcript_id", event.TranscriptID))
	return nil
}

func (p *Producer) PublishProcessingFailed(ctx context.Context, event ProcessingFailedEvent) error {
	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyProcessingFailed, event)
	if err != nil {
		return fmt.Errorf("failed to publish processing failed event: %w", err)
	}

	slog.Error("Published processing failed event",
		slog.String("video_id", event.VideoID),
		slog.String("stage", event.Stage),
		slog.String("error", event.Error))
	return nil
}

func (p *Producer) PublishVideoReceived(ctx context.Context, event VideoReceivedEvent) error {
	err := p.rabbitmq.Publish(ctx, ExchangeTopic, RoutingKeyVideoReceived, event)
	if err != nil {
		return fmt.Errorf("failed to publish video received event: %w", err)
	}

	slog.Info("Published video received event",
		slog.String("video_id", event.VideoID),
		slog.String("url", event.URL))
	return nil
}
