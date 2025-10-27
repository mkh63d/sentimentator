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
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/transcript"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

type Server struct {
	port              int
	db                database.Service
	transcriptHandler handlers.TranscriptHandler
}

func NewServer() *http.Server {
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
	videoService := video.NewService()
	transcriptService := transcript.NewService(videoService, whisperAPIKey)
	transcriptHandler := handlers.NewTranscriptHandler(transcriptService)

	NewServer := &Server{
		port:              port,
		db:                database.New(),
		transcriptHandler: *transcriptHandler,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
