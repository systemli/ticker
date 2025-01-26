package response

import (
	"time"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Ticker struct {
	ID          int         `json:"id"`
	CreatedAt   time.Time   `json:"createdAt"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Active      bool        `json:"active"`
	Information Information `json:"information"`
	Websites    []Website   `json:"websites"`
	Telegram    Telegram    `json:"telegram"`
	Mastodon    Mastodon    `json:"mastodon"`
	Bluesky     Bluesky     `json:"bluesky"`
	SignalGroup SignalGroup `json:"signalGroup"`
	Location    Location    `json:"location"`
}

type Information struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
	Telegram string `json:"telegram"`
	Mastodon string `json:"mastodon"`
	Bluesky  string `json:"bluesky"`
}

type Website struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Origin    string    `json:"origin"`
}

type Telegram struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	BotUsername string `json:"botUsername"`
	ChannelName string `json:"channelName"`
}

type Mastodon struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	Name        string `json:"name"`
	Server      string `json:"server"`
	ScreenName  string `json:"screenName"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type Bluesky struct {
	Active    bool   `json:"active"`
	Connected bool   `json:"connected"`
	Handle    string `json:"handle"`
}

type SignalGroup struct {
	Active          bool   `json:"active"`
	Connected       bool   `json:"connected"`
	GroupID         string `json:"groupID"`
	GroupInviteLink string `json:"groupInviteLink"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func TickerResponse(t storage.Ticker, config config.Config) Ticker {
	websites := make([]Website, 0)
	for _, website := range t.Websites {
		websites = append(websites, Website{
			ID:        website.ID,
			CreatedAt: website.CreatedAt,
			Origin:    website.Origin,
		})
	}

	return Ticker{
		ID:          t.ID,
		CreatedAt:   t.CreatedAt,
		Title:       t.Title,
		Description: t.Description,
		Active:      t.Active,
		Information: Information{
			Author:   t.Information.Author,
			URL:      t.Information.URL,
			Email:    t.Information.Email,
			Twitter:  t.Information.Twitter,
			Facebook: t.Information.Facebook,
			Telegram: t.Information.Telegram,
			Mastodon: t.Information.Mastodon,
			Bluesky:  t.Information.Bluesky,
		},
		Websites: websites,
		Telegram: Telegram{
			Active:      t.Telegram.Active,
			Connected:   t.Telegram.Connected(),
			BotUsername: config.Telegram.User.UserName,
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
		Bluesky: Bluesky{
			Active:    t.Bluesky.Active,
			Connected: t.Bluesky.Connected(),
			Handle:    t.Bluesky.Handle,
		},
		SignalGroup: SignalGroup{
			Active:          t.SignalGroup.Active,
			Connected:       t.SignalGroup.Connected(),
			GroupID:         t.SignalGroup.GroupID,
			GroupInviteLink: t.SignalGroup.GroupInviteLink,
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
