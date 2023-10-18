package response

import (
	"fmt"
	"time"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Timeline []TimelineEntry

type TimelineEntry struct {
	ID             int                 `json:"id"`
	CreatedAt      time.Time           `json:"createdAt"`
	Text           string              `json:"text"`
	GeoInformation string              `json:"geoInformation"`
	Attachments    []MessageAttachment `json:"attachments"`
}

type Attachment struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

func TimelineResponse(messages []storage.Message, config config.Config) []TimelineEntry {
	timeline := make([]TimelineEntry, 0)
	for _, message := range messages {
		m, _ := message.GeoInformation.MarshalJSON()

		var attachments []MessageAttachment
		for _, attachment := range message.Attachments {
			name := fmt.Sprintf("%s.%s", attachment.UUID, attachment.Extension)
			attachments = append(attachments, MessageAttachment{URL: MediaURL(config.UploadURL, name), ContentType: attachment.ContentType})
		}

		timeline = append(timeline, TimelineEntry{
			ID:             message.ID,
			CreatedAt:      message.CreatedAt,
			Text:           message.Text,
			GeoInformation: string(m),
			Attachments:    attachments,
		})

	}
	return timeline
}
