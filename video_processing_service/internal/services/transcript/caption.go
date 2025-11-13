package transcript

import (
	"context"
	"strings"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	"github.com/Oskarowski/sentimentator/video_processing_service/internal/services/extractor"
)

type CaptionStrategy struct {
	extractor extractor.Service
}

func NewCaptionStrategy(extractorService extractor.Service) *CaptionStrategy {
	return &CaptionStrategy{
		extractor: extractorService,
	}
}

func (c *CaptionStrategy) Name() string {
	return "captions"
}

func (c *CaptionStrategy) GetTranscript(
	ctx context.Context,
	videoURL, platformName, videoID, lang string,
) (*Result, error) {
	// Use extractor service to retrieve captions directly
	segments, err := c.extractor.GetCaptions(ctx, videoURL, lang)
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

	// Convert extractor segments to models segments
	mSegments := make([]models.Segment, len(segments))
	for i, seg := range segments {
		mSegments[i] = models.Segment{Start: seg.Start, End: seg.End, Text: seg.Text}
	}

	return &Result{
		Text:     strings.TrimSpace(textBuilder.String()),
		Segments: mSegments,
		Method:   models.TranscriptMethod(c.Name()),
	}, nil

}
