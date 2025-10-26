package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"io/fs"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	handlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/api"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/health", handlers.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/status", handlers.ServiceStatus)
	}

	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic("failed to find 'assets' subdirectory in embedded filesystem: " + err.Error())
	}
	r.StaticFS("/assets", http.FS(staticFiles))

	r.GET("/video", func(c *gin.Context) {
		component := web.VideoUrlForm()
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/video/process", func(c *gin.Context) {
		web.VideoUrlWebHandler(c.Writer, c.Request)
	})

	return r
}
