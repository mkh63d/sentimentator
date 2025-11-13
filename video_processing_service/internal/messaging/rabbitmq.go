package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

type Config struct {
	URL           string
	MaxRetries    int
	RetryDelay    time.Duration
	PrefetchCount int
	PrefetchSize  int
	ConsumerTag   string
}

func NewRabbitMQ(config Config) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var err error

	maxRetries := config.MaxRetries
	if maxRetries == 0 {
		maxRetries = 5
	}

	retryDelay := config.RetryDelay
	if retryDelay == 0 {
		retryDelay = 5 * time.Second
	}

	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(config.URL)
		if err == nil {
			break
		}
		slog.Warn("Failed to connect to RabbitMQ, retrying...",
			slog.Int("attempt", i+1),
			slog.Int("max_retries", maxRetries),
			slog.String("error", err.Error()))
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS if configured
	if config.PrefetchCount > 0 {
		err = channel.Qos(
			config.PrefetchCount,
			config.PrefetchSize,
			false, // global
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	slog.Info("Successfully connected to RabbitMQ")

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		url:     config.URL,
	}, nil
}

func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) DeclareQueueWithArgs(name string, args map[string]interface{}) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
}

func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue binds a queue to an exchange with a routing key
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

func (r *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	slog.Debug("Message published",
		slog.String("exchange", exchange),
		slog.String("routing_key", routingKey))

	return nil
}

// Consume starts consuming messages from a queue
func (r *RabbitMQ) Consume(queueName, consumerTag string) (<-chan amqp.Delivery, error) {
	if consumerTag == "" {
		consumerTag = "video-processing-service"
	}

	msgs, err := r.channel.Consume(
		queueName,   // queue
		consumerTag, // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	slog.Info("Started consuming messages",
		slog.String("queue", queueName),
		slog.String("consumer_tag", consumerTag))

	return msgs, nil
}

// Close closes the channel and connection
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			slog.Error("Failed to close channel", slog.String("error", err.Error()))
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			slog.Error("Failed to close connection", slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("RabbitMQ connection closed")
	return nil
}

func (r *RabbitMQ) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}

func (r *RabbitMQ) Reconnect(config Config) error {
	r.Close()

	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel on reconnect: %w", err)
	}

	if config.PrefetchCount > 0 {
		err = channel.Qos(config.PrefetchCount, config.PrefetchSize, false)
		if err != nil {
			channel.Close()
			conn.Close()
			return fmt.Errorf("failed to set QoS on reconnect: %w", err)
		}
	}

	r.conn = conn
	r.channel = channel

	slog.Info("Successfully reconnected to RabbitMQ")
	return nil
}
