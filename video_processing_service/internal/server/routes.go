package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"io/fs"

	"github.com/Oskarowski/sentimentator/video_processing_service/cmd/web"
	handlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/api"
	webhandlers "github.com/Oskarowski/sentimentator/video_processing_service/internal/handlers/web"

	"github.com/a-h/templ"
)

// authMiddleware protects /web routes using a static bearer token
// The token is provided via env var WEB_ADMIN_TOKEN. Requests must include
// either header Authorization: Bearer <token> or a cookie named "web_admin_token"
// with the same value.
func authMiddleware() gin.HandlerFunc {
	expected := strings.TrimSpace(os.Getenv("WEB_ADMIN_TOKEN"))
	return func(c *gin.Context) {
		// If no token is configured, allow access
		if expected == "" {
			c.Next()
			return
		}

		auth := c.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(auth, "Bearer "); ok {
			token := strings.TrimSpace(after)
			if token == expected {
				c.Next()
				return
			}
		}

		if cookie, err := c.Cookie("web_admin_token"); err == nil {
			if strings.TrimSpace(cookie) == expected {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://localhost:8000",
		"https://sentimentator.jdex.pl",
		"https://sentimentator.api.jdex.pl",
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * 3600,
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

	videoWebHandler := webhandlers.NewVideoHandler(s.db, s.extractorService, s.transcriptService, s.producer)
	debugHandler := webhandlers.NewDebugHandler(s.producer)

	r.StaticFS("/web/assets", http.FS(staticFiles))

	r.GET("/web/login", func(c *gin.Context) {
		component := web.LoginPage()
		templ.Handler(component).ServeHTTP(c.Writer, c.Request)
	})

	wr := r.Group("/web")
	wr.Use(authMiddleware())
	{
		r.POST("/web/login", func(c *gin.Context) {
			token := c.PostForm("token")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "token required"})
				return
			}
			c.SetCookie("web_admin_token", token, 12*60*60, "/", "", false, true)
			c.Redirect(http.StatusFound, "/web/")
		})

		wr.GET("/", func(c *gin.Context) {
			component := web.MainPage()
			templ.Handler(component).ServeHTTP(c.Writer, c.Request)
		})

		wr.GET("/video", func(c *gin.Context) {
			component := web.MainPage()
			templ.Handler(component).ServeHTTP(c.Writer, c.Request)
		})

		wr.GET("/videos/list", videoWebHandler.ListVideos)

		wr.GET("/debug/events", debugHandler.ServeDebugPage)
		wr.POST("/debug/events/video-received", debugHandler.PublishVideoReceived)
		wr.POST("/debug/events/transcript-ready", debugHandler.PublishTranscriptReady)
		wr.POST("/debug/events/processing-failed", debugHandler.PublishProcessingFailed)
	}

	return r
}
