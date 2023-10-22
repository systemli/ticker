package response

import (
	"fmt"
	"time"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Message struct {
	ID             int                 `json:"id"`
	CreatedAt      time.Time           `json:"createdAt"`
	Text           string              `json:"text"`
	Ticker         int                 `json:"ticker"`
	TelegramURL    string              `json:"telegramUrl,omitempty"`
	MastodonURL    string              `json:"mastodonUrl,omitempty"`
	GeoInformation string              `json:"geoInformation"`
	Attachments    []MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
}

func MessageResponse(message storage.Message, config config.Config) Message {
	m, _ := message.GeoInformation.MarshalJSON()
	var attachments []MessageAttachment

	for _, attachment := range message.Attachments {
		name := fmt.Sprintf("%s.%s", attachment.UUID, attachment.Extension)
		attachments = append(attachments, MessageAttachment{URL: MediaURL(config.UploadURL, name), ContentType: attachment.ContentType})
	}

	return Message{
		ID:             message.ID,
		CreatedAt:      message.CreatedAt,
		Text:           message.Text,
		Ticker:         message.TickerID,
		TelegramURL:    message.TelegramURL(),
		MastodonURL:    message.MastodonURL(),
		GeoInformation: string(m),
		Attachments:    attachments,
	}
}

func MessagesResponse(messages []storage.Message, config config.Config) []Message {
	msgs := make([]Message, 0)
	for _, message := range messages {
		msgs = append(msgs, MessageResponse(message, config))
	}
	return msgs
}

func MediaURL(uploadURL, name string) string {
	return fmt.Sprintf("%s/media/%s", uploadURL, name)
}
