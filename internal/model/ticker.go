package model

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

//Ticker represents the structure of an Ticker configuration
type Ticker struct {
	ID           int       `storm:"id,increment"`
	CreationDate time.Time `storm:"index"`
	Domain       string    `storm:"unique"`
	Title        string
	Description  string
	Active       bool
	Hashtags     []string
	Information  Information
	Twitter      Twitter
	Location     Location
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

//Location represents the default location for the Ticker.
type Location struct {
	Lat float64
	Lon float64
}

//TickerResponse represents the Ticker for API responses.
type TickerResponse struct {
	ID           int                 `json:"id"`
	CreationDate time.Time           `json:"creation_date"`
	Domain       string              `json:"domain"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Active       bool                `json:"active"`
	Hashtags     []string            `json:"hashtags"`
	Information  InformationResponse `json:"information"`
	Twitter      TwitterResponse     `json:"twitter"`
	Location     LocationResponse    `json:"location"`
}

//InformationResponse represents the Information for API responses.
type InformationResponse struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
}

//TwitterResponse represents the Twitter settings for API responses.
type TwitterResponse struct {
	Active      bool   `json:"active"`
	Connected   bool   `json:"connected"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

//LocationResponse represents the Location for API responses.
type LocationResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

//NewTicker creates new Ticker
func NewTicker() *Ticker {
	return &Ticker{
		CreationDate: time.Now(),
	}
}

//Reset set most variables to there defaults
func (t *Ticker) Reset() {
	t.Active = false
	t.Description = ""
	t.Hashtags = []string{}
	t.Information = Information{}
	t.Twitter.Secret = ""
	t.Twitter.Token = ""
	t.Twitter.Active = false
	t.Twitter.User = twitter.User{}
	t.Location = Location{}
}

//NewTickerResponse returns a API friendly representation for a Ticker.
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

	l := LocationResponse{
		Lat: ticker.Location.Lat,
		Lon: ticker.Location.Lon,
	}

	return &TickerResponse{
		ID:           ticker.ID,
		CreationDate: ticker.CreationDate,
		Domain:       ticker.Domain,
		Title:        ticker.Title,
		Description:  ticker.Description,
		Active:       ticker.Active,
		Hashtags:     ticker.Hashtags,
		Information:  info,
		Twitter:      tw,
		Location:     l,
	}
}

//NewTickersResponse prepares a map of []TickerResponse.
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
