package response

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Message struct {
	ID               int                 `json:"id"`
	CreationDate     time.Time           `json:"creation_date"`
	Text             string              `json:"text"`
	Ticker           int                 `json:"ticker"`
	TweetID          string              `json:"tweet_id"`
	TweetUser        string              `json:"tweet_user"`
	TelegramMessages []tgbotapi.Message  `json:"telegram_messages"`
	GeoInformation   string              `json:"geo_information"`
	Attachments      []MessageAttachment `json:"attachments"`
}

type MessageAttachment struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

func MessageResponse(message storage.Message, config config.Config) Message {
	m, _ := message.GeoInformation.MarshalJSON()
	var attachments []MessageAttachment

	for _, attachment := range message.Attachments {
		name := fmt.Sprintf("%s.%s", attachment.UUID, attachment.Extension)
		attachments = append(attachments, MessageAttachment{URL: MediaURL(config.UploadURL, name), ContentType: attachment.ContentType})
	}

	return Message{
		ID:               message.ID,
		CreationDate:     message.CreationDate,
		Text:             message.Text,
		Ticker:           message.Ticker,
		TweetID:          message.Tweet.ID,
		TweetUser:        message.Tweet.UserName,
		TelegramMessages: message.Telegram.Messages,
		GeoInformation:   string(m),
		Attachments:      attachments,
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
