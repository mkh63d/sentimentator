package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"io/fs"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	handlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/api"
	webhandlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/web"

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
		v1.POST("/transcript", s.transcriptHandler.GetTranscript)

	}

	staticFiles, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic("failed to find 'assets' subdirectory in embedded filesystem: " + err.Error())
	}

	videoWebHandler := webhandlers.NewVideoHandler(s.db, s.videoService, s.transcriptService, s.producer)
	debugHandler := webhandlers.NewDebugHandler(s.producer)

	wr := r.Group("/web")
	{
		wr.StaticFS("/assets", http.FS(staticFiles))
		wr.GET("/", func(c *gin.Context) {
			component := web.MainPage()
			templ.Handler(component).ServeHTTP(c.Writer, c.Request)
		})

		wr.GET("/video", func(c *gin.Context) {
			component := web.MainPage()
			templ.Handler(component).ServeHTTP(c.Writer, c.Request)
		})

		wr.POST("/video/process", videoWebHandler.ProcessVideo)

		wr.GET("/jobs/stats", videoWebHandler.GetJobStats)

		wr.GET("/videos/list", videoWebHandler.ListVideos)

		wr.GET("/debug/events", debugHandler.ServeDebugPage)
		wr.POST("/debug/events/video-received", debugHandler.PublishVideoReceived)
		wr.POST("/debug/events/transcript-ready", debugHandler.PublishTranscriptReady)
		wr.POST("/debug/events/processing-failed", debugHandler.PublishProcessingFailed)
	}

	return r
}
