package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(ctx context.Context, body []byte) error

type Consumer struct {
	rabbitmq *RabbitMQ
	handlers map[string]MessageHandler
}

func NewConsumer(rabbitmq *RabbitMQ) *Consumer {
	return &Consumer{
		rabbitmq: rabbitmq,
		handlers: make(map[string]MessageHandler),
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

	err := handler(ctx, msg.Body)
	if err != nil {
		slog.Error("Failed to process message",
			slog.String("queue", queueName),
			slog.String("error", err.Error()))

		// Reject and requeue the message (or send to DLQ based on your config)
		if err := msg.Nack(false, true); err != nil {
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

func HandleVideoReceivedEvent(ctx context.Context, body []byte) error {
	var event VideoReceivedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	slog.Info("Received video event",
		slog.String("video_id", event.VideoID),
		slog.String("url", event.URL),
		slog.String("platform", event.Platform))

	// TODO: Implement actual video processing logic here

	return nil
}
