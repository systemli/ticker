package response

import (
	"time"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Ticker struct {
	ID           int         `json:"id"`
	CreationDate time.Time   `json:"creation_date"`
	Domain       string      `json:"domain"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Active       bool        `json:"active"`
	Information  Information `json:"information"`
	Telegram     Telegram    `json:"telegram"`
	Mastodon     Mastodon    `json:"mastodon"`
	Location     Location    `json:"location"`
}

type Information struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
	Telegram string `json:"telegram"`
}

type Telegram struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	BotUsername string `json:"bot_username"`
	ChannelName string `json:"channel_name"`
}

type Mastodon struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	Name        string `json:"name"`
	Server      string `json:"server"`
	ScreenName  string `json:"screen_name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func TickerResponse(t storage.Ticker, config config.Config) Ticker {
	return Ticker{
		ID:           t.ID,
		CreationDate: t.CreatedAt,
		Domain:       t.Domain,
		Title:        t.Title,
		Description:  t.Description,
		Active:       t.Active,
		Information: Information{
			Author:   t.Information.Author,
			URL:      t.Information.URL,
			Email:    t.Information.Email,
			Twitter:  t.Information.Twitter,
			Facebook: t.Information.Facebook,
			Telegram: t.Information.Telegram,
		},
		Telegram: Telegram{
			Active:      t.Telegram.Active,
			Connected:   t.Telegram.Connected(),
			BotUsername: config.TelegramBotUser.UserName,
			ChannelName: t.Telegram.ChannelName,
		},
		Mastodon: Mastodon{
			Active:     t.Mastodon.Active,
			Connected:  t.Mastodon.Connected(),
			Name:       t.Mastodon.User.Username,
			Server:     t.Mastodon.Server,
			ScreenName: t.Mastodon.User.DisplayName,
			ImageURL:   t.Mastodon.User.Avatar,
		},
		Location: Location{
			Lat: t.Location.Lat,
			Lon: t.Location.Lon,
		},
	}
}

func TickersResponse(tickers []storage.Ticker, config config.Config) []Ticker {
	t := make([]Ticker, 0)

	for _, ticker := range tickers {
		t = append(t, TickerResponse(ticker, config))
	}
	return t
}
