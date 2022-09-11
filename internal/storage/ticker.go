package storage

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type Ticker struct {
	ID           int       `storm:"id,increment"`
	CreationDate time.Time `storm:"index"`
	Domain       string    `storm:"unique"`
	Title        string
	Description  string
	Active       bool
	Information  Information
	Twitter      Twitter
	Telegram     Telegram
	Location     Location
}

func NewTicker() Ticker {
	return Ticker{
		CreationDate: time.Now(),
	}
}

func (t *Ticker) Reset() {
	t.Active = false
	t.Description = ""
	t.Information = Information{}
	t.Twitter.Secret = ""
	t.Twitter.Token = ""
	t.Twitter.Active = false
	t.Twitter.User = twitter.User{}
	t.Telegram.Active = false
	t.Telegram.ChannelName = ""
	t.Location = Location{}
}

type Information struct {
	Author   string
	URL      string
	Email    string
	Twitter  string
	Facebook string
	Telegram string
}

type Twitter struct {
	Active bool
	Token  string
	Secret string
	User   twitter.User
}

func (tw *Twitter) Connected() bool {
	return tw.Token != "" && tw.Secret != ""
}

type Telegram struct {
	Active      bool   `json:"active"`
	ChannelName string `json:"channel_name"`
}

type Location struct {
	Lat float64
	Lon float64
}
