package storage

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	geojson "github.com/paulmach/go.geojson"
)

type Message struct {
	ID             int `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TickerID       int `gorm:"index"`
	Text           string
	Attachments    []Attachment
	GeoInformation geojson.FeatureCollection `gorm:"serializer:json"`
	Telegram       TelegramMeta              `gorm:"serializer:json"`
	Mastodon       MastodonMeta              `gorm:"serializer:json"`
}

func NewMessage() Message {
	return Message{}
}

type TelegramMeta struct {
	Messages []tgbotapi.Message
}

type MastodonMeta struct {
	ID  string
	URI string
	URL string
}

type Attachment struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MessageID   int `gorm:"index"`
	UUID        string
	Extension   string
	ContentType string
}

func (m *Message) AddAttachment(upload Upload) {
	attachment := Attachment{
		UUID:        upload.UUID,
		Extension:   upload.Extension,
		ContentType: upload.ContentType,
	}

	m.Attachments = append(m.Attachments, attachment)
}

func (m *Message) AddAttachments(uploads []Upload) {
	for _, upload := range uploads {
		m.AddAttachment(upload)
	}
}

func (m *Message) TelegramURL() string {
	if len(m.Telegram.Messages) == 0 {
		return ""
	}

	message := m.Telegram.Messages[0]
	return fmt.Sprintf("https://t.me/%s/%d", message.Chat.UserName, message.MessageID)
}

func (m *Message) MastodonURL() string {
	if m.Mastodon.ID == "" {
		return ""
	}

	return m.Mastodon.URL
}
