package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TickerID    int `gorm:"index"`
	Text        string
	Attachments []Attachment
	Telegram    TelegramMeta    `gorm:"serializer:json"`
	Mastodon    MastodonMeta    `gorm:"serializer:json"`
	Bluesky     BlueskyMeta     `gorm:"serializer:json"`
	SignalGroup SignalGroupMeta `gorm:"serializer:json"`
	MatrixRoom  MatrixRoomMeta  `gorm:"serializer:json"`
}

func NewMessage() Message {
	return Message{}
}

func (m *Message) AsMap() map[string]interface{} {
	telegram, _ := json.Marshal(m.Telegram)
	mastodon, _ := json.Marshal(m.Mastodon)
	bluesky, _ := json.Marshal(m.Bluesky)
	signalGroup, _ := json.Marshal(m.SignalGroup)
	matrixRoom, _ := json.Marshal(m.MatrixRoom)

	return map[string]interface{}{
		"id":           m.ID,
		"created_at":   m.CreatedAt,
		"updated_at":   m.UpdatedAt,
		"ticker_id":    m.TickerID,
		"text":         m.Text,
		"telegram":     telegram,
		"mastodon":     mastodon,
		"bluesky":      bluesky,
		"signal_group": signalGroup,
		"matrix_room":  matrixRoom,
	}
}

type TelegramMeta struct {
	Messages []tgbotapi.Message
}

type MastodonMeta struct {
	ID  string
	URI string
	URL string
}

type BlueskyMeta struct {
	Handle string
	Uri    string
	Cid    string
}

type SignalGroupMeta struct {
	Timestamp int
}

type MatrixRoomMeta struct {
	EventIDs []string
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

func (m *Message) BlueskyURL() string {
	parts := strings.Split(m.Bluesky.Uri, "/")
	if len(parts) < 3 {
		return ""
	}

	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", m.Bluesky.Handle, parts[len(parts)-1])
}
