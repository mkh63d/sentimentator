package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/database"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/utils"
)

func VideoUrlWebHandler(w http.ResponseWriter, r *http.Request, db database.Service) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ytUrl := r.FormValue("yt_url_input")
	if ytUrl == "" {
		http.Error(w, "YouTube URL is required", http.StatusBadRequest)
		return
	}

	// Parse the video URL to extract video ID, platform, and original URL
	videoInfo, err := utils.ParseVideoURL(ytUrl)
	if err != nil {
		component := VideoUrlError(err.Error())
		renderErr := component.Render(r.Context(), w)
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
			log.Printf("Error rendering error component: %v", renderErr)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	videoID, err := db.SaveVideo(ctx, videoInfo.VideoID, videoInfo.Platform, videoInfo.OriginalURL)
	if err != nil {
		log.Printf("Error saving video to database: %v", err)
		component := VideoUrlError("Failed to save video to database")
		renderErr := component.Render(r.Context(), w)
		if renderErr != nil {
			http.Error(w, renderErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	component := VideoUrlPost(videoInfo.VideoID, videoInfo.Platform, videoInfo.OriginalURL, videoID)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error rendering in VideoUrlWebHandler: %v", err)
	}
}
