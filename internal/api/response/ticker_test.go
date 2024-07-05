package response

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type TickersResponseTestSuite struct {
	suite.Suite
}

func (s *TickersResponseTestSuite) TestTickersResponse() {
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
			Mastodon: "https://systemli.social/@example",
			Bluesky:  "https://example.com",
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
		SignalGroup: storage.TickerSignalGroup{
			Active:          true,
			GroupID:         "example",
			GroupInviteLink: "https://signal.group/#example",
		},
		Location: storage.TickerLocation{
			Lat: 0.0,
			Lon: 0.0,
		},
	}

	config := config.Config{
		Telegram: config.Telegram{
			User: tgbotapi.User{
				UserName: "ticker",
			},
		},
	}

	tickerResponse := TickersResponse([]storage.Ticker{ticker}, config)

	s.Equal(1, len(tickerResponse))
	s.Equal(ticker.ID, tickerResponse[0].ID)
	s.Equal(ticker.CreatedAt, tickerResponse[0].CreatedAt)
	s.Equal(ticker.Domain, tickerResponse[0].Domain)
	s.Equal(ticker.Title, tickerResponse[0].Title)
	s.Equal(ticker.Description, tickerResponse[0].Description)
	s.Equal(ticker.Active, tickerResponse[0].Active)
	s.Equal(ticker.Information.Author, tickerResponse[0].Information.Author)
	s.Equal(ticker.Information.URL, tickerResponse[0].Information.URL)
	s.Equal(ticker.Information.Email, tickerResponse[0].Information.Email)
	s.Equal(ticker.Information.Twitter, tickerResponse[0].Information.Twitter)
	s.Equal(ticker.Information.Facebook, tickerResponse[0].Information.Facebook)
	s.Equal(ticker.Information.Telegram, tickerResponse[0].Information.Telegram)
	s.Equal(ticker.Information.Mastodon, tickerResponse[0].Information.Mastodon)
	s.Equal(ticker.Information.Bluesky, tickerResponse[0].Information.Bluesky)
	s.Equal(ticker.Telegram.Active, tickerResponse[0].Telegram.Active)
	s.Equal(ticker.Telegram.Connected(), tickerResponse[0].Telegram.Connected)
	s.Equal(config.Telegram.User.UserName, tickerResponse[0].Telegram.BotUsername)
	s.Equal(ticker.Telegram.ChannelName, tickerResponse[0].Telegram.ChannelName)
	s.Equal(ticker.Mastodon.Active, tickerResponse[0].Mastodon.Active)
	s.Equal(ticker.Mastodon.Connected(), tickerResponse[0].Mastodon.Connected)
	s.Equal(ticker.Mastodon.User.Username, tickerResponse[0].Mastodon.Name)
	s.Equal(ticker.Mastodon.Server, tickerResponse[0].Mastodon.Server)
	s.Equal(ticker.Mastodon.User.DisplayName, tickerResponse[0].Mastodon.ScreenName)
	s.Equal(ticker.Mastodon.User.Avatar, tickerResponse[0].Mastodon.ImageURL)
	s.Equal(ticker.SignalGroup.Active, tickerResponse[0].SignalGroup.Active)
	s.Equal(ticker.SignalGroup.Connected(), tickerResponse[0].SignalGroup.Connected)
	s.Equal(ticker.SignalGroup.GroupID, tickerResponse[0].SignalGroup.GroupID)
	s.Equal(ticker.Location.Lat, tickerResponse[0].Location.Lat)
	s.Equal(ticker.Location.Lon, tickerResponse[0].Location.Lon)
}

func TestTickersResponseTestSuite(t *testing.T) {
	suite.Run(t, new(TickersResponseTestSuite))
}
