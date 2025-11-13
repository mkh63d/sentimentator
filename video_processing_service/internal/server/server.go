package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	handlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/api"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/messaging"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
)

type Server struct {
	port              int
	db                database.Service
	extractorService  extractor.Service
	transcriptService *transcript.Service
	transcriptHandler handlers.TranscriptHandler
	rabbitmq          *messaging.RabbitMQ
	producer          *messaging.Producer
	consumer          *messaging.Consumer
	consumerCtx       context.Context
	consumerCancel    context.CancelFunc
}

func NewServer() (*Server, *http.Server) {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		portEnv = "8080"
	}
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		port = 8080
	}

	rabbitmqURL := os.Getenv("AMQP_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
		slog.Warn("AMQP_URL not set, using default", slog.String("url", rabbitmqURL))
	}

	db := database.New()
	if err := database.AutoMigrate(db); err != nil {
		panic(fmt.Sprintf("failed to run database migrations: %v", err))
	}

	rabbitmq, err := messaging.NewRabbitMQ(messaging.Config{
		URL:           rabbitmqURL,
		MaxRetries:    2,
		RetryDelay:    5 * time.Second,
		PrefetchCount: 5,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
	}

	if err := setupRabbitMQTopology(rabbitmq); err != nil {
		panic(fmt.Sprintf("failed to setup RabbitMQ topology: %v", err))
	}

	// Initialize services
	extractorService := extractor.NewService("")
	transcriptService := transcript.NewService(extractorService)

	producer := messaging.NewProducer(rabbitmq)
	consumer := messaging.NewConsumer(rabbitmq, db, transcriptService, extractorService, producer)
	transcriptHandler := handlers.NewTranscriptHandler(db, extractorService, producer)

	consumer.RegisterHandler(messaging.QueueVideoProcessing, consumer.HandleTranscriptRequestEvent)

	consumerCtx, consumerCancel := context.WithCancel(context.Background())

	srv := &Server{
		port:              port,
		db:                db,
		extractorService:  extractorService,
		transcriptService: transcriptService,
		transcriptHandler: *transcriptHandler,
		rabbitmq:          rabbitmq,
		producer:          producer,
		consumer:          consumer,
		consumerCtx:       consumerCtx,
		consumerCancel:    consumerCancel,
	}

	if err := consumer.StartConsuming(consumerCtx); err != nil {
		panic(fmt.Sprintf("failed to start consuming messages: %v", err))
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv, httpServer
}

// setupRabbitMQTopology declares exchanges, queues, and bindings
func setupRabbitMQTopology(rabbitmq *messaging.RabbitMQ) error {
	if err := rabbitmq.DeclareExchange(messaging.ExchangeTopic, "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	if err := rabbitmq.DeclareExchange("sentimentator.dlx", "topic"); err != nil {
		return fmt.Errorf("failed to declare dead letter exchange: %w", err)
	}

	queues := []string{
		messaging.QueueVideoProcessing,
		messaging.QueueTranscript,
	}

	for _, queueName := range queues {
		args := map[string]interface{}{
			"x-dead-letter-exchange":    "sentimentator.dlx",
			"x-dead-letter-routing-key": queueName + ".failed",
			"x-message-ttl":             int64(15 * time.Minute / time.Millisecond),
		}
		queue, err := rabbitmq.DeclareQueueWithArgs(queueName, args)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
		slog.Info("Declared queue with DLX",
			slog.String("queue", queueName),
			slog.Int("messages", queue.Messages))
	}

	for _, queueName := range queues {
		dlqName := queueName + ".failed"
		if _, err := rabbitmq.DeclareQueue(dlqName); err != nil {
			return fmt.Errorf("failed to declare dead letter queue %s: %w", dlqName, err)
		}

		if err := rabbitmq.BindQueue(dlqName, queueName+".failed", "sentimentator.dlx"); err != nil {
			return fmt.Errorf("failed to bind dead letter queue %s: %w", dlqName, err)
		}
	}

	bindings := map[string]string{
		messaging.QueueVideoProcessing: messaging.RoutingKeyVideoTranscriptRequested,
		messaging.QueueTranscript:      messaging.RoutingKeyVideoTranscriptCompleted,
	}

	for queueName, routingKey := range bindings {
		if err := rabbitmq.BindQueue(queueName, routingKey, messaging.ExchangeTopic); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	slog.Info("RabbitMQ topology setup completed")
	return nil
}

// Shutdown gracefully stops the consumer and closes RabbitMQ connection
func (s *Server) Shutdown() {
	slog.Info("Shutting down server...")

	if s.consumerCancel != nil {
		s.consumerCancel()
	}

	if s.rabbitmq != nil {
		if err := s.rabbitmq.Close(); err != nil {
			slog.Error("Failed to close RabbitMQ connection", slog.String("error", err.Error()))
		}
	}

	slog.Info("Server shutdown complete")
}
