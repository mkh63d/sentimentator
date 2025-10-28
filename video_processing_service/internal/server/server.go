package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	handlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/api"
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

	db := database.New()
	if err := database.AutoMigrate(db); err != nil {
		panic(fmt.Sprintf("failed to run database migrations: %v", err))
	}

	videoService := video.NewService()
	transcriptService := transcript.NewService(videoService, whisperAPIKey)
	queueService := queue.NewService(db)
	transcriptHandler := handlers.NewTranscriptHandler(db, videoService, queueService)

	workerPool := queue.NewWorkerPool(db, transcriptService, 1, 5*time.Second)

	srv := &Server{
		port:              port,
		db:                db,
		videoService:      videoService,
		transcriptService: transcriptService,
		queueService:      queueService,
		transcriptHandler: *transcriptHandler,
		workerPool:        workerPool,
	}

	workerPool.Start()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv, httpServer
}

// Shutdown gracefully stops the worker pool
func (s *Server) Shutdown() {
	if s.workerPool != nil {
		s.workerPool.Stop()
	}
}
