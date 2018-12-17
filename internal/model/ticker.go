package model

import (
	"github.com/dghubble/go-twitter/twitter"
	"time"
)

//Ticker represents the structure of an Ticker configuration
type Ticker struct {
	ID           int       `storm:"id,increment"`
	CreationDate time.Time `storm:"index"`
	Domain       string    `storm:"unique"`
	Title        string
	Description  string
	Active       bool
	PrependTime  bool `json:"prepend_time"`
	Hashtags     []string
	Information  Information
	Twitter      Twitter
}

//Information holds some meta information for Ticker
type Information struct {
	Author   string
	URL      string
	Email    string
	Twitter  string
	Facebook string
}

//Twitter holds all required twitter information.
type Twitter struct {
	Active bool
	Token  string
	Secret string
	User   twitter.User
}

type TickerResponse struct {
	ID           int                 `json:"id"`
	CreationDate time.Time           `json:"creation_date"`
	Domain       string              `json:"domain"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Active       bool                `json:"active"`
	PrependTime  bool                `json:"prepend_time"`
	Hashtags     []string            `json:"hashtags"`
	Information  InformationResponse `json:"information"`
	Twitter      TwitterResponse     `json:"twitter"`
}

type InformationResponse struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
}

type TwitterResponse struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

//NewTicker creates new Ticker
func NewTicker() *Ticker {
	return &Ticker{
		CreationDate: time.Now(),
	}
}

//
func NewTickerResponse(ticker *Ticker) *TickerResponse {
	info := InformationResponse{
		Author:   ticker.Information.Author,
		URL:      ticker.Information.URL,
		Email:    ticker.Information.Email,
		Twitter:  ticker.Information.Twitter,
		Facebook: ticker.Information.Facebook,
	}

	tw := TwitterResponse{
		Active:      ticker.Twitter.Active,
		Connected:   ticker.Twitter.Connected(),
		Name:        ticker.Twitter.User.Name,
		ScreenName:  ticker.Twitter.User.ScreenName,
		Description: ticker.Twitter.User.Description,
		ImageURL:    ticker.Twitter.User.ProfileImageURLHttps,
	}

	return &TickerResponse{
		ID:           ticker.ID,
		CreationDate: ticker.CreationDate,
		Domain:       ticker.Domain,
		Title:        ticker.Title,
		Description:  ticker.Description,
		Active:       ticker.Active,
		PrependTime:  ticker.PrependTime,
		Hashtags:     ticker.Hashtags,
		Information:  info,
		Twitter:      tw,
	}
}

//
func NewTickersResponse(tickers []*Ticker) []*TickerResponse {
	var tr []*TickerResponse

	for _, ticker := range tickers {
		tr = append(tr, NewTickerResponse(ticker))
	}

	return tr
}

//Connected returns true when twitter can be used.
func (tw *Twitter) Connected() bool {
	return tw.Token != "" && tw.Secret != ""
}
