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
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/queue"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

type Server struct {
	port              int
	db                database.Service
	videoService      *video.Service
	transcriptService *transcript.Service
	queueService      *queue.Service
	transcriptHandler handlers.TranscriptHandler
	workerPool        *queue.WorkerPool
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

	whisperAPIKey := os.Getenv("OPENAI_API_KEY")
	if whisperAPIKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
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
		MaxRetries:    5,
		RetryDelay:    5 * time.Second,
		PrefetchCount: 1,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
	}

	if err := setupRabbitMQTopology(rabbitmq); err != nil {
		panic(fmt.Sprintf("failed to setup RabbitMQ topology: %v", err))
	}

	// Initialize services
	videoService := video.NewService()
	transcriptService := transcript.NewService(videoService, whisperAPIKey)
	queueService := queue.NewService(db)
	transcriptHandler := handlers.NewTranscriptHandler(db, videoService, queueService)

	workerPool := queue.NewWorkerPool(db, transcriptService, 1, 5*time.Second)

	producer := messaging.NewProducer(rabbitmq)
	consumer := messaging.NewConsumer(rabbitmq, db, transcriptService, videoService, producer)

	consumer.RegisterHandler(messaging.QueueVideoProcessing, consumer.HandleVideoReceivedEvent)

	consumerCtx, consumerCancel := context.WithCancel(context.Background())

	srv := &Server{
		port:              port,
		db:                db,
		videoService:      videoService,
		transcriptService: transcriptService,
		queueService:      queueService,
		transcriptHandler: *transcriptHandler,
		workerPool:        workerPool,
		rabbitmq:          rabbitmq,
		producer:          producer,
		consumer:          consumer,
		consumerCtx:       consumerCtx,
		consumerCancel:    consumerCancel,
	}

	workerPool.Start()

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
	// Declare topic exchange
	if err := rabbitmq.DeclareExchange(messaging.ExchangeTopic, "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queues := []string{
		messaging.QueueVideoProcessing,
		messaging.QueueTranscript,
		messaging.QueueDeadLetter,
	}

	for _, queueName := range queues {
		if _, err := rabbitmq.DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	// Bind queues to exchange with routing keys
	bindings := map[string]string{
		messaging.QueueVideoProcessing: messaging.RoutingKeyVideoReceived,
		messaging.QueueTranscript:      messaging.RoutingKeyTranscriptReady,
	}

	for queueName, routingKey := range bindings {
		if err := rabbitmq.BindQueue(queueName, routingKey, messaging.ExchangeTopic); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	slog.Info("RabbitMQ topology setup completed")
	return nil
}

// Shutdown gracefully stops the worker pool, consumer, and closes RabbitMQ connection
func (s *Server) Shutdown() {
	slog.Info("Shutting down server...")

	if s.consumerCancel != nil {
		s.consumerCancel()
	}

	if s.workerPool != nil {
		s.workerPool.Stop()
	}

	if s.rabbitmq != nil {
		if err := s.rabbitmq.Close(); err != nil {
			slog.Error("Failed to close RabbitMQ connection", slog.String("error", err.Error()))
		}
	}

	slog.Info("Server shutdown complete")
}
