package legacy

import (
	"time"

	"github.com/mattn/go-mastodon"
)

type Ticker struct {
	ID           int       `storm:"id,increment"`
	CreationDate time.Time `storm:"index"`
	Domain       string    `storm:"unique"`
	Title        string
	Description  string
	Active       bool
	Information  Information
	Telegram     Telegram
	Mastodon     Mastodon
	Location     Location
}

type Information struct {
	Author   string
	URL      string
	Email    string
	Twitter  string
	Facebook string
	Telegram string
}

type Telegram struct {
	Active      bool   `json:"active"`
	ChannelName string `json:"channel_name"`
}

type Mastodon struct {
	Active      bool   `json:"active"`
	Server      string `json:"server"`
	Token       string `json:"token"`
	Secret      string `json:"secret"`
	AccessToken string `json:"access_token"`
	User        mastodon.Account
}

type Location struct {
	Lat float64
	Lon float64
}
