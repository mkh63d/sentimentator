package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(ctx context.Context, body []byte) error

type Consumer struct {
	rabbitmq          *RabbitMQ
	handlers          map[string]MessageHandler
	db                database.Service
	extractorService  extractor.Service
	transcriptService *transcript.Service
	producer          *Producer
}

func NewConsumer(rabbitmq *RabbitMQ, db database.Service, transcriptService *transcript.Service, extractorService extractor.Service, producer *Producer) *Consumer {
	return &Consumer{
		rabbitmq:          rabbitmq,
		handlers:          make(map[string]MessageHandler),
		db:                db,
		extractorService:  extractorService,
		transcriptService: transcriptService,
		producer:          producer,
	}
}

func (c *Consumer) RegisterHandler(queueName string, handler MessageHandler) {
	c.handlers[queueName] = handler
	slog.Info("Registered handler", slog.String("queue", queueName))
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	for queueName, handler := range c.handlers {
		consumerTag := fmt.Sprintf("vps-%s", queueName)

		msgs, err := c.rabbitmq.channel.Consume(
			queueName,
			consumerTag,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to start consuming from queue %s: %w", queueName, err)
		}

		slog.Info("Started consuming messages",
			slog.String("queue", queueName),
			slog.String("consumer_tag", consumerTag))

		go c.processMessages(ctx, queueName, msgs, handler)
	}

	return nil
}

func (c *Consumer) processMessages(ctx context.Context, queueName string, msgs <-chan amqp.Delivery, handler MessageHandler) {
	slog.Info("Starting message processor", slog.String("queue", queueName))

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping message processor", slog.String("queue", queueName))
			return
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Message channel closed", slog.String("queue", queueName))
				return
			}

			c.handleMessage(ctx, msg, handler, queueName)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg amqp.Delivery, handler MessageHandler, queueName string) {
	slog.Debug("Processing message",
		slog.String("queue", queueName),
		slog.String("message_id", msg.MessageId),
		slog.String("correlation_id", msg.CorrelationId))

	processCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	err := handler(processCtx, msg.Body)
	if err != nil {
		slog.Error("Failed to process message",
			slog.String("queue", queueName),
			slog.String("error", err.Error()))

		if processCtx.Err() == context.DeadlineExceeded {
			slog.Error("Message processing timeout exceeded",
				slog.String("queue", queueName),
				slog.String("timeout", "15 minutes"))
		}

		requeue := true
		if xDeath, ok := msg.Headers["x-death"].([]interface{}); ok && len(xDeath) > 0 {
			if deathMap, ok := xDeath[0].(amqp.Table); ok {
				if count, ok := deathMap["count"].(int64); ok && count >= 3 {
					slog.Warn("Message exceeded retry limit, rejecting permanently",
						slog.String("queue", queueName),
						slog.Int64("retry_count", count))
					requeue = false
				}
			}
		}

		if err := msg.Nack(false, requeue); err != nil {
			slog.Error("Failed to nack message", slog.String("error", err.Error()))
		}
		return
	}

	// Acknowledge the message
	if err := msg.Ack(false); err != nil {
		slog.Error("Failed to ack message", slog.String("error", err.Error()))
	}

	slog.Debug("Message processed successfully", slog.String("queue", queueName))
}

func (c *Consumer) HandleTranscriptRequestEvent(ctx context.Context, body []byte) error {
	startTime := time.Now()

	var event VideoTranscriptRequestedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	slog.Info("Received transcript request event",
		slog.String("url", event.URL),
		slog.String("analysis_id", event.AnalysisID))

	metadata, err := c.extractorService.GetMetadata(ctx, event.URL)
	if err != nil {
		return fmt.Errorf("failed to get video metadata: %w", err)
	}

	slog.Info("Retrieved video metadata",
		slog.String("video_id", metadata.VideoID),
		slog.String("title", metadata.Title),
		slog.Int("duration", metadata.Duration),
		slog.String("elapsed", time.Since(startTime).String()))

	videoRecord, err := c.db.GetVideoByVideoID(ctx, metadata.VideoID, metadata.Platform)
	if err != nil {
		videoRecord = &models.Video{
			VideoID:     metadata.VideoID,
			Platform:    metadata.Platform,
			OriginalURL: event.URL,
			Title:       metadata.Title,
			Duration:    metadata.Duration,
		}

		if err := c.db.SaveVideo(ctx, videoRecord); err != nil {
			return fmt.Errorf("failed to save video record: %w", err)
		}
		slog.Info("Created new video record", slog.String("video_id", metadata.VideoID))
	} else {
		slog.Info("Found existing video record", slog.String("video_id", metadata.VideoID))
	}

	hasTranscript, existingTranscript, err := c.db.HasTranscript(ctx, videoRecord.ID, "en")
	if err != nil {
		return fmt.Errorf("failed to check for existing transcript: %w", err)
	}

	if hasTranscript {
		if existingTranscript.Status == models.StatusCompleted {
			transcriptEvent := VideoTranscriptCompletedEvent{
				AnalysisID:   event.AnalysisID,
				VideoID:      videoRecord.VideoID,
				TranscriptID: fmt.Sprintf("%d", existingTranscript.ID),
				Text:         existingTranscript.Text,
				Language:     existingTranscript.Language,
				Duration:     videoRecord.Duration,
				CreatedAt:    time.Now(),
			}

			err := c.producer.PublishTranscriptReady(ctx, transcriptEvent)
			if err != nil {
				return fmt.Errorf("failed to publish transcript ready event: %w", err)
			}

			slog.Info("Transcript already exists, published ready event",
				slog.String("video_id", videoRecord.VideoID),
				slog.String("transcript_id", fmt.Sprintf("%d", existingTranscript.ID)),
				slog.String("total_elapsed", time.Since(startTime).String()),
			)
			return nil
		}
	}

	slog.Info("Starting transcript generation",
		slog.String("video_id", metadata.VideoID),
		slog.String("url", event.URL))

	transcriptStartTime := time.Now()
	result, err := c.transcriptService.GetTranscript(ctx, event.URL, "en")
	if err != nil {
		slog.Error("Failed to get transcript",
			slog.String("video_id", videoRecord.VideoID),
			slog.String("url", event.URL),
			slog.String("error", err.Error()),
			slog.String("elapsed", time.Since(transcriptStartTime).String()),
		)
		return fmt.Errorf("failed to get transcript: %w", err)
	}

	slog.Info("Transcript generated successfully",
		slog.String("video_id", metadata.VideoID),
		slog.Int("text_length", len(result.Text)),
		slog.String("method", string(result.Method)),
		slog.String("transcript_elapsed", time.Since(transcriptStartTime).String()))

	transcriptRecord := &models.Transcript{
		VideoID:   videoRecord.ID,
		Language:  "en",
		Method:    result.Method,
		Status:    models.StatusCompleted,
		Text:      result.Text,
		Segments:  result.Segments,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	if err := c.db.CreateTranscript(ctx, transcriptRecord); err != nil {
		return fmt.Errorf("failed to save transcript record: %w", err)
	}

	transcriptEvent := VideoTranscriptCompletedEvent{
		AnalysisID:   event.AnalysisID,
		VideoID:      videoRecord.VideoID,
		TranscriptID: fmt.Sprintf("%d", transcriptRecord.ID),
		Text:         result.Text,
		Language:     transcriptRecord.Language,
		Duration:     videoRecord.Duration,
		CreatedAt:    time.Now(),
	}

	if err := c.producer.PublishTranscriptReady(ctx, transcriptEvent); err != nil {
		slog.Error("Failed to publish transcript completed event",
			slog.String("video_id", videoRecord.VideoID),
			slog.String("transcript_id", fmt.Sprintf("%d", transcriptRecord.ID)),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to publish transcript ready event: %w", err)
	}

	slog.Info("Transcript processing completed successfully",
		slog.String("video_id", metadata.VideoID),
		slog.String("transcript_id", fmt.Sprintf("%d", transcriptRecord.ID)),
		slog.String("total_elapsed", time.Since(startTime).String()))

	return nil
}
