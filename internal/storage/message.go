package storage

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mattn/go-mastodon"
	geojson "github.com/paulmach/go.geojson"
)

type Message struct {
	ID             int       `storm:"id,increment"`
	CreationDate   time.Time `storm:"index"`
	Ticker         int       `storm:"index"`
	Text           string
	Attachments    []Attachment
	GeoInformation geojson.FeatureCollection
	Tweet          Tweet
	Telegram       TelegramMeta
	Mastodon       mastodon.Status
}

func NewMessage() Message {
	return Message{
		CreationDate: time.Now(),
	}
}

type Tweet struct {
	ID       string
	UserName string
}

type TelegramMeta struct {
	Messages []tgbotapi.Message
}

type Attachment struct {
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
