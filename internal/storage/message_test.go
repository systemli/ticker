package storage

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mattn/go-mastodon"
	"github.com/stretchr/testify/assert"
)

func TestAddAttachments(t *testing.T) {
	upload := NewUpload("image.jpg", "image/jped", 1)
	message := NewMessage()
	message.AddAttachments([]Upload{upload})

	assert.Equal(t, 1, len(message.Attachments))
}

func TestTwitterURL(t *testing.T) {
	message := NewMessage()

	assert.Empty(t, message.TwitterURL())

	message.Tweet.ID = "1"
	message.Tweet.UserName = "systemli"

	assert.Equal(t, "https://twitter.com/systemli/status/1", message.TwitterURL())
}

func TestTelegramURL(t *testing.T) {
	message := NewMessage()

	assert.Empty(t, message.TelegramURL())

	message.Telegram = TelegramMeta{
		Messages: []tgbotapi.Message{
			{
				MessageID: 1,
				Chat: &tgbotapi.Chat{
					UserName: "systemli",
				}},
		},
	}

	assert.Equal(t, "https://t.me/systemli/1", message.TelegramURL())

}

func TestMastodonURL(t *testing.T) {
	message := NewMessage()

	assert.Empty(t, message.MastodonURL())

	url := "https://mastodon.social/web/@systemli/1"
	message.Mastodon = mastodon.Status{
		ID:  "1",
		URL: url,
	}

	assert.Equal(t, url, message.MastodonURL())
}
