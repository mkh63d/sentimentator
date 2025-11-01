package transcript

import (
	"context"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/video"
)

type CaptionStrategy struct {
	videoService *video.Service
}

func NewCaptionStrategy(videoService *video.Service) *CaptionStrategy {
	return &CaptionStrategy{
		videoService: videoService,
	}
}

func (c *CaptionStrategy) Name() string {
	return "captions"
}

func (c *CaptionStrategy) GetTranscript(
	ctx context.Context,
	videoURL, platformName, videoID, lang string,
) (*Result, error) {
	platform, err := c.videoService.DetectPlatform(videoURL)
	if err != nil {
		return nil, err
	}

	// Check if platform supports captions before attempting to fetch
	segments, err := platform.GetCaptions(ctx, videoID, lang)
	if err != nil {
		if strings.Contains(err.Error(), "does not provide captions") {
			return nil, err
		}
		return nil, err
	}

	if len(segments) == 0 {
		return nil, models.ErrNoCaptionsAvailable
	}

	var textBuilder strings.Builder
	for _, seg := range segments {
		textBuilder.WriteString(seg.Text)
		textBuilder.WriteString(" ")
	}

	return &Result{
		Text:     strings.TrimSpace(textBuilder.String()),
		Segments: segments,
		Method:   models.TranscriptMethod(c.Name()),
	}, nil

}
