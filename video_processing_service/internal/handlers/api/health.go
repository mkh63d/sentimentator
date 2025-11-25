package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

type ServiceStatusResponse struct {
	Status       string            `json:"status"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Timestamp    time.Time         `json:"timestamp"`
	Dependencies map[string]string `json:"dependencies"`
}

func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Service:   "video_processing_service",
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

func ServiceStatus(c *gin.Context) {
	response := ServiceStatusResponse{
		Status:    "healthy",
		Service:   "video_processing_service",
		Version:   "0.1.0",
		Timestamp: time.Now(),
		Dependencies: map[string]string{
			"database":           "health_status_check_not_configured",
			"rabbitmq":           "health_status_check_not_configured",
			"openai_whisper_api": "health_status_check_not_configured",
			"yt-dlp":             "health_status_check_not_configured",
		},
	}

	c.JSON(http.StatusOK, response)
}
