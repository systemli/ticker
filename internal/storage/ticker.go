package storage

import (
	"net/http"
	"strconv"
	"time"
)

type Ticker struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Domain      string `gorm:"unique,index"`
	Title       string
	Description string
	Active      bool
	Information TickerInformation `gorm:"embedded"`
	Location    TickerLocation    `gorm:"embedded"`
	Telegram    TickerTelegram
	Mastodon    TickerMastodon
	Users       []User `gorm:"many2many:ticker_users;"`
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

func (t *Ticker) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"created_at":  t.CreatedAt,
		"updated_at":  t.UpdatedAt,
		"domain":      t.Domain,
		"title":       t.Title,
		"description": t.Description,
		"active":      t.Active,
		"author":      t.Information.Author,
		"url":         t.Information.URL,
		"email":       t.Information.Email,
		"twitter":     t.Information.Twitter,
		"facebook":    t.Information.Facebook,
		"telegram":    t.Information.Telegram,
		"lat":         t.Location.Lat,
		"lon":         t.Location.Lon,
	}
}

type TickerInformation struct {
	Author   string
	URL      string
	Email    string
	Twitter  string
	Facebook string
	Telegram string
	Mastodon string
	Bluesky  string
}

type TickerTelegram struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TickerID    int `gorm:"index"`
	Active      bool
	ChannelName string
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
	TickerID    int `gorm:"index"`
	Active      bool
	Server      string
	Token       string
	Secret      string
	AccessToken string
	User        MastodonUser `gorm:"embedded"`
}

type MastodonUser struct {
	Username    string
	DisplayName string
	Avatar      string
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
	m.User = MastodonUser{}
}

type TickerLocation struct {
	Lat float64
	Lon float64
}

type TickerFilter struct {
	Domain  *string
	Title   *string
	Active  *bool
	OrderBy string
	Sort    string
}

func NewTickerFilter(req *http.Request) TickerFilter {
	filter := TickerFilter{
		OrderBy: "id",
		Sort:    "asc",
	}

	if req == nil {
		return filter
	}

	if req.URL.Query().Get("order_by") != "" {
		opts := []string{"id", "created_at", "updated_at", "domain", "title", "active"}
		for _, opt := range opts {
			if req.URL.Query().Get("order_by") == opt {
				filter.OrderBy = req.URL.Query().Get("order_by")
				break
			}
		}
	}

	if req.URL.Query().Get("sort") == "asc" {
		filter.Sort = "asc"
	} else {
		filter.Sort = "desc"
	}

	domain := req.URL.Query().Get("domain")
	if domain != "" {
		filter.Domain = &domain
	}

	title := req.URL.Query().Get("title")
	if title != "" {
		filter.Title = &title
	}

	active := req.URL.Query().Get("active")
	if active != "" {
		activeBool, _ := strconv.ParseBool(active)
		filter.Active = &activeBool
	}

	return filter
}
