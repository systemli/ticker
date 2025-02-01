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
	Domain      string
	Title       string
	Description string
	Active      bool
	Information TickerInformation `gorm:"embedded"`
	Location    TickerLocation    `gorm:"embedded"`
	Telegram    TickerTelegram
	Mastodon    TickerMastodon
	Bluesky     TickerBluesky
	SignalGroup TickerSignalGroup
	Websites    []TickerWebsite `gorm:"foreignKey:TickerID;"`
	Users       []User          `gorm:"many2many:ticker_users;"`
}

func NewTicker() Ticker {
	return Ticker{}
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
		"threads":     t.Information.Threads,
		"instagram":   t.Information.Instagram,
		"telegram":    t.Information.Telegram,
		"bluesky":     t.Information.Bluesky,
		"mastodon":    t.Information.Mastodon,
		"lat":         t.Location.Lat,
		"lon":         t.Location.Lon,
	}
}

type TickerInformation struct {
	Author    string
	URL       string
	Email     string
	Twitter   string
	Facebook  string
	Instagram string
	Threads   string
	Telegram  string
	Mastodon  string
	Bluesky   string
}

type TickerWebsite struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TickerID  int    `gorm:"index;not null"`
	Origin    string `gorm:"unique;not null"`
}

type TickerTelegram struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TickerID    int `gorm:"index"`
	Active      bool
	ChannelName string
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

type TickerBluesky struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TickerID  int `gorm:"index"`
	Active    bool
	Handle    string
	// AppKey is the application password from Bluesky
	// Future consideration: persist the access token, refresh token instead of app key
	AppKey string
}

func (b *TickerBluesky) Connected() bool {
	return b.Handle != "" && b.AppKey != ""
}

type TickerSignalGroup struct {
	ID              int `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	TickerID        int `gorm:"index"`
	Active          bool
	GroupID         string
	GroupInviteLink string
}

func (s *TickerSignalGroup) Connected() bool {
	return s.GroupID != ""
}

type TickerLocation struct {
	Lat float64
	Lon float64
}

type TickerFilter struct {
	Origin  *string
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
		opts := []string{"id", "created_at", "updated_at", "origin", "title", "active"}
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

	origin := req.URL.Query().Get("origin")
	if origin != "" {
		filter.Origin = &origin
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
