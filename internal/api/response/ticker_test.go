package response

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func TestTickersResponse(t *testing.T) {
	ticker := storage.Ticker{
		ID:          1,
		CreatedAt:   time.Now(),
		Domain:      "example.com",
		Title:       "Example",
		Description: "Example",
		Active:      true,
		Information: storage.TickerInformation{
			Author:   "Example",
			URL:      "https://example.com",
			Email:    "contact@example.com",
			Twitter:  "@example",
			Facebook: "https://facebook.com/example",
			Telegram: "example",
		},
		Telegram: storage.TickerTelegram{
			Active:      true,
			ChannelName: "example",
		},
		Mastodon: storage.TickerMastodon{
			Active: true,
			Server: "https://example.com",
			User: storage.MastodonUser{
				Username:    "example",
				DisplayName: "Example",
				Avatar:      "https://example.com/avatar.png",
			},
		},
		Location: storage.TickerLocation{
			Lat: 0.0,
			Lon: 0.0,
		},
	}

	config := config.Config{
		TelegramBotUser: tgbotapi.User{
			UserName: "ticker",
		},
	}

	tickerResponse := TickersResponse([]storage.Ticker{ticker}, config)

	assert.Equal(t, 1, len(tickerResponse))
	assert.Equal(t, ticker.ID, tickerResponse[0].ID)
	assert.Equal(t, ticker.CreatedAt, tickerResponse[0].CreationDate)
	assert.Equal(t, ticker.Domain, tickerResponse[0].Domain)
	assert.Equal(t, ticker.Title, tickerResponse[0].Title)
	assert.Equal(t, ticker.Description, tickerResponse[0].Description)
	assert.Equal(t, ticker.Active, tickerResponse[0].Active)
	assert.Equal(t, ticker.Information.Author, tickerResponse[0].Information.Author)
	assert.Equal(t, ticker.Information.URL, tickerResponse[0].Information.URL)
	assert.Equal(t, ticker.Information.Email, tickerResponse[0].Information.Email)
	assert.Equal(t, ticker.Information.Twitter, tickerResponse[0].Information.Twitter)
	assert.Equal(t, ticker.Information.Facebook, tickerResponse[0].Information.Facebook)
	assert.Equal(t, ticker.Information.Telegram, tickerResponse[0].Information.Telegram)
	assert.Equal(t, ticker.Telegram.Active, tickerResponse[0].Telegram.Active)
	assert.Equal(t, ticker.Telegram.Connected(), tickerResponse[0].Telegram.Connected)
	assert.Equal(t, config.TelegramBotUser.UserName, tickerResponse[0].Telegram.BotUsername)
	assert.Equal(t, ticker.Telegram.ChannelName, tickerResponse[0].Telegram.ChannelName)
	assert.Equal(t, ticker.Mastodon.Active, tickerResponse[0].Mastodon.Active)
	assert.Equal(t, ticker.Mastodon.Connected(), tickerResponse[0].Mastodon.Connected)
	assert.Equal(t, ticker.Mastodon.User.Username, tickerResponse[0].Mastodon.Name)
	assert.Equal(t, ticker.Mastodon.Server, tickerResponse[0].Mastodon.Server)
	assert.Equal(t, ticker.Mastodon.User.DisplayName, tickerResponse[0].Mastodon.ScreenName)
	assert.Equal(t, ticker.Mastodon.User.Avatar, tickerResponse[0].Mastodon.ImageURL)
	assert.Equal(t, ticker.Location.Lat, tickerResponse[0].Location.Lat)
	assert.Equal(t, ticker.Location.Lon, tickerResponse[0].Location.Lon)
}
