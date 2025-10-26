package web

import (
	"log"
	"net/http"
)

func VideoUrlWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	ytUrl := r.FormValue("yt_url_input")
	if ytUrl == "" {
		http.Error(w, "YouTube URL is required", http.StatusBadRequest)
		return
	}
	component := VideoUrlPost(ytUrl)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
