package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type Message struct {
	ID          int                 `json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	Text        string              `json:"text"`
	Ticker      int                 `json:"ticker"`
	TelegramURL string              `json:"telegramUrl,omitempty"`
	MastodonURL string              `json:"mastodonUrl,omitempty"`
	BlueskyURL  string              `json:"blueskyUrl,omitempty"`
	Attachments []MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
}

func MessageResponse(message storage.Message) Message {
	var attachments []MessageAttachment

	for _, attachment := range message.Attachments {
		attachments = append(attachments, MessageAttachment{URL: storage.MediaURL(attachment.FileName()), ContentType: attachment.ContentType})
	}

	return Message{
		ID:          message.ID,
		CreatedAt:   message.CreatedAt,
		Text:        message.Text,
		Ticker:      message.TickerID,
		TelegramURL: message.TelegramURL(),
		MastodonURL: message.MastodonURL(),
		BlueskyURL:  message.BlueskyURL(),
		Attachments: attachments,
	}
}

func MessagesResponse(messages []storage.Message) []Message {
	msgs := make([]Message, 0)
	for _, message := range messages {
		msgs = append(msgs, MessageResponse(message))
	}
	return msgs
}
