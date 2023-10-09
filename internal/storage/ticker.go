package storage

import (
	"time"

	"github.com/mattn/go-mastodon"
	"gorm.io/gorm"
)

type Ticker struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Domain      string `gorm:"unique"`
	Title       string
	Description string
	Active      bool
	Information TickerInformation
	Telegram    TickerTelegram
	Mastodon    TickerMastodon
	Location    TickerLocation
	Users       []User `gorm:"many2many:user_tickers;"`
}

func NewTicker() Ticker {
	return Ticker{}
}

func (t *Ticker) Reset() {
	t.Active = false
	t.Description = ""
	t.Information = TickerInformation{}
	t.Location = TickerLocation{}

	t.Telegram.Reset()
	t.Mastodon.Reset()
}

type TickerInformation struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TickerID  int
	Author    string
	URL       string
	Email     string
	Twitter   string
	Facebook  string
	Telegram  string
}

type TickerTelegram struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TickerID    int    `gorm:"index"`
	Active      bool   `json:"active"`
	ChannelName string `json:"channel_name"`
}

func (tg *TickerTelegram) Reset() {
	tg.Active = false
	tg.ChannelName = ""
}

func (tg *TickerTelegram) Connected() bool {
	return tg.ChannelName != ""
}

type TickerMastodon struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TickerID    int              `gorm:"index"`
	Active      bool             `json:"active"`
	Server      string           `json:"server"`
	Token       string           `json:"token"`
	Secret      string           `json:"secret"`
	AccessToken string           `json:"access_token"`
	User        mastodon.Account `gorm:"type:json"`
}

func (m *TickerMastodon) Connected() bool {
	return m.Token != "" && m.Secret != "" && m.AccessToken != ""
}

func (m *TickerMastodon) Reset() {
	m.Active = false
	m.Server = ""
	m.Token = ""
	m.Secret = ""
	m.AccessToken = ""
	m.User = mastodon.Account{}
}

type TickerLocation struct {
	gorm.Model
	ID       int `gorm:"primaryKey"`
	TickerID int
	Lat      float64
	Lon      float64
}
