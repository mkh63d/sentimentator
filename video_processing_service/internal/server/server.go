package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
)

type Server struct {
	port int

	db database.Service
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

	NewServer := &Server{
		port: port,

		db: database.New(),
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
