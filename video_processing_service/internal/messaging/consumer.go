package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"strconv"
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
	workerPool        chan struct{} // Semaphore Pattern https://blog.stackademic.com/how-to-implement-a-semaphore-concurrency-pattern-in-golang-3a51e5e517a8
}

func NewConsumer(rabbitmq *RabbitMQ, db database.Service, transcriptService *transcript.Service, extractorService extractor.Service, producer *Producer) *Consumer {
	maxConcurrentWorkersEnv, exists := os.LookupEnv("MAX_CONCURRENT_RABBIT_PROCESSING_WORKERS")
	if !exists {
		maxConcurrentWorkersEnv = "5"
		slog.Warn("MAX_CONCURRENT_RABBIT_PROCESSING_WORKERS not set, using default", slog.String("max_workers", maxConcurrentWorkersEnv))
	}
	maxConcurrentWorkers, err := strconv.Atoi(maxConcurrentWorkersEnv)
	if err != nil || maxConcurrentWorkers <= 0 {
		maxConcurrentWorkers = 5
		slog.Warn("Invalid MAX_CONCURRENT_RABBIT_PROCESSING_WORKERS value, using default", slog.Int("max_workers", maxConcurrentWorkers))
	}

	return &Consumer{
		rabbitmq:          rabbitmq,
		handlers:          make(map[string]MessageHandler),
		db:                db,
		extractorService:  extractorService,
		transcriptService: transcriptService,
		producer:          producer,
		workerPool:        make(chan struct{}, maxConcurrentWorkers),
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
	retryCount := int64(0)
	if msg.Headers != nil {
		if val, ok := msg.Headers["retry-count"]; ok {
			switch v := val.(type) {
			case int64:
				retryCount = v
			case int32:
				retryCount = int64(v)
			case int:
				retryCount = int64(v)
			case float64:
				retryCount = int64(v)
			case string:
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					retryCount = parsed
				}
			}
		}
	}

	// async processing with worker pool
	if queueName == QueueVideoProcessing {
		select {
		case c.workerPool <- struct{}{}:
			// this is for getting a worker slot for processing the message
		case <-ctx.Done():
			slog.Warn("Context cancelled while waiting for worker slot, requeuing message",
				slog.String("queue", queueName))
			if err := msg.Nack(false, true); err != nil {
				slog.Error("Failed to nack message", slog.String("error", err.Error()))
			}
			return
		}

		// Acknowledge immediately to allow processing next message
		if err := msg.Ack(false); err != nil {
			slog.Error("Failed to ack message", slog.String("error", err.Error()))
			<-c.workerPool
			return
		}

		go func() {
			defer func() { <-c.workerPool }()

			processCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
			defer cancel()

			err := handler(processCtx, msg.Body)
			if err != nil {
				slog.Error("Failed to process message asynchronously",
					slog.String("queue", queueName),
					slog.String("error", err.Error()),
					slog.Int64("retry_count", retryCount))

				if processCtx.Err() == context.DeadlineExceeded {
					slog.Error("Message processing timeout exceeded",
						slog.String("queue", queueName),
						slog.String("timeout", "15 minutes"))
				}

				if retryCount < 2 {
					slog.Info("Re-publishing message for retry",
						slog.String("queue", queueName),
						slog.Int64("retry_count", retryCount))

					if err := c.republishForRetry(context.Background(), msg, queueName, retryCount+1); err != nil {
						slog.Error("Failed to republish message for retry",
							slog.String("error", err.Error()))
					}
				} else {
					slog.Warn("Message exceeded retry limit, discarding permanently",
						slog.String("queue", queueName),
						slog.Int64("retry_count", retryCount))
				}
			} else {
				slog.Debug("Message processed successfully (async)", slog.String("queue", queueName))
			}
		}()
		return
	}

	// sync processing for other queues
	processCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	err := handler(processCtx, msg.Body)

	if err != nil {
		slog.Error("Failed to process message",
			slog.String("queue", queueName),
			slog.String("error", err.Error()),
			slog.Int64("retry_count", retryCount))

		if processCtx.Err() == context.DeadlineExceeded {
			slog.Error("Message processing timeout exceeded",
				slog.String("queue", queueName),
				slog.String("timeout", "15 minutes"))
		}

		if retryCount >= 2 {
			slog.Warn("Message exceeded retry limit, discarding",
				slog.String("queue", queueName),
				slog.Int64("retry_count", retryCount))
			// Reject without requeue - message is discarded permanently
			if err := msg.Reject(false); err != nil {
				slog.Error("Failed to reject message", slog.String("error", err.Error()))
			}
			return
		}

		if err := msg.Nack(false, false); err != nil {
			slog.Error("Failed to nack message", slog.String("error", err.Error()))
		}

		slog.Info("Re-publishing message for retry",
			slog.String("queue", queueName),
			slog.Int64("retry_count", retryCount))

		if err := c.republishForRetry(context.Background(), msg, queueName, retryCount+1); err != nil {
			slog.Error("Failed to republish message for retry",
				slog.String("queue", queueName),
				slog.String("error", err.Error()))
		}
		return
	}

	// Acknowledge the message on success
	if err := msg.Ack(false); err != nil {
		slog.Error("Failed to ack message",
			slog.String("queue", queueName),
			slog.String("error", err.Error()))
	}
}

// republishForRetry publishes a failed message back to the queue with retry metadata
func (c *Consumer) republishForRetry(ctx context.Context, originalMsg amqp.Delivery, queueName string, retryCount int64) error {
	headers := make(map[string]any)
	if originalMsg.Headers != nil {
		maps.Copy(headers, originalMsg.Headers)
	}

	headers["retry-count"] = retryCount

	return c.rabbitmq.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   originalMsg.ContentType,
			Body:          originalMsg.Body,
			Timestamp:     time.Now(),
			MessageId:     originalMsg.MessageId,
			CorrelationId: originalMsg.CorrelationId,
			Headers:       headers,
		},
	)
}

func (c *Consumer) HandleTranscriptRequestEvent(ctx context.Context, body []byte) error {
	startTime := time.Now()

	var event VideoTranscriptRequestedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	slog.Info("Processing transcript request",
		slog.String("analysis_id", event.AnalysisID),
		slog.String("url", event.URL))

	metadata, err := c.extractorService.GetMetadata(ctx, event.URL)
	if err != nil {
		return fmt.Errorf("failed to get video metadata: %w", err)
	}

	slog.Info("Video metadata retrieved",
		slog.String("video_id", metadata.VideoID),
		slog.String("title", metadata.Title))

	videoRecord, err := c.db.GetVideoByVideoID(ctx, metadata.VideoID, metadata.Platform)
	if err != nil {
		videoRecord = &models.Video{
			VideoID:      metadata.VideoID,
			Platform:     metadata.Platform,
			OriginalURL:  event.URL,
			Title:        metadata.Title,
			Duration:     metadata.Duration,
			ThumbnailURL: metadata.ThumbnailURL,
		}

		if err := c.db.SaveVideo(ctx, videoRecord); err != nil {
			return fmt.Errorf("failed to save video record: %w", err)
		}
		slog.Debug("Created new video record", slog.String("video_id", metadata.VideoID))
	} else {
		slog.Debug("Found existing video record", slog.String("video_id", metadata.VideoID))
	}

	hasTranscript, existingTranscript, err := c.db.HasTranscript(ctx, videoRecord)
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
			slog.Info("Transcript already exists, reusing",
				slog.String("video_id", videoRecord.VideoID),
				slog.String("transcript_id", fmt.Sprintf("%d", existingTranscript.ID)))
			return nil
		} else {
			slog.Warn("Transcript in progress, skipping",
				slog.String("video_id", videoRecord.VideoID),
				slog.String("status", string(existingTranscript.Status)))
			return fmt.Errorf("transcript already exists with status: %s", existingTranscript.Status)
		}
	}

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

	slog.Debug("Transcript generated",
		slog.String("video_id", metadata.VideoID),
		slog.String("method", string(result.Method)),
		slog.Int64("total_processing_seconds", result.TotalProcessingDuration))

	transcriptRecord := &models.Transcript{
		VideoID:                 videoRecord.ID,
		Language:                result.DetectedLanguage,
		Method:                  result.Method,
		Status:                  models.StatusCompleted,
		Text:                    result.Text,
		Segments:                result.Segments,
		AudioExtractionDuration: &result.AudioExtractionDuration,
		TranscriptionDuration:   &result.TranscriptionDuration,
		TotalProcessingDuration: &result.TotalProcessingDuration,
		AudioFileSizeBytes:      &result.AudioFileSizeBytes,
		UpdatedAt:               time.Now(),
		CreatedAt:               time.Now(),
	}
	if err := c.db.CreateTranscript(ctx, transcriptRecord); err != nil {
		return fmt.Errorf("failed to save transcript record: %w", err)
	}

	transcriptEvent := VideoTranscriptCompletedEvent{
		AnalysisID:   event.AnalysisID,
		OriginalURL:  videoRecord.OriginalURL,
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

	slog.Info("Transcript retrieval completed",
		slog.String("analysis_id", event.AnalysisID),
		slog.String("video_id", metadata.VideoID),
		slog.String("transcript_id", fmt.Sprintf("%d", transcriptRecord.ID)),
		slog.Int64("total_seconds", int64(time.Since(startTime).Seconds())))

	return nil
}
