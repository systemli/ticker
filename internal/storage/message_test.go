package storage

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestAddAttachments(t *testing.T) {
	upload := NewUpload("image.jpg", "image/jped", 1)
	message := NewMessage()
	message.AddAttachments([]Upload{upload})

	assert.Equal(t, 1, len(message.Attachments))
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
	message.Mastodon = MastodonMeta{
		ID:  "1",
		URL: url,
	}

	assert.Equal(t, url, message.MastodonURL())
}

func TestBlueskyURL(t *testing.T) {
	message := NewMessage()

	assert.Empty(t, message.BlueskyURL())

	message.Bluesky = BlueskyMeta{
		Uri: "",
	}

	assert.Empty(t, message.BlueskyURL())

	url := "https://bsky.app/profile/systemli.bsky.social/post/3kr7p3jxkpw2n"
	message.Bluesky = BlueskyMeta{
		Uri:    "at://did:plc:izpk4tc54wu6b3yufcdixqje/app.bsky.feed.post/3kr7p3jxkpw2n",
		Handle: "systemli.bsky.social",
	}

	assert.Equal(t, url, message.BlueskyURL())
}
