package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type Timeline []TimelineEntry

type TimelineEntry struct {
	ID          int          `json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
}

func TimelineResponse(messages []storage.Message) []TimelineEntry {
	timeline := make([]TimelineEntry, 0)
	for _, message := range messages {
		var attachments []Attachment
		for _, attachment := range message.Attachments {
			attachments = append(attachments, Attachment{URL: storage.MediaURL(attachment.FileName()), ContentType: attachment.ContentType})
		}

		timeline = append(timeline, TimelineEntry{
			ID:          message.ID,
			CreatedAt:   message.CreatedAt,
			Text:        message.Text,
			Attachments: attachments,
		})

	}
	return timeline
}
