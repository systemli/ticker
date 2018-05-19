package model

import "time"

//Ticker represents the structure of an Ticker configuration
type Ticker struct {
	ID           int         `json:"id" storm:"id,increment"`
	CreationDate time.Time   `json:"creation_date" storm:"index"`
	Domain       string      `json:"domain" binding:"required" storm:"unique"`
	Title        string      `json:"title" binding:"required"`
	Description  string      `json:"description" binding:"required"`
	Active       bool        `json:"active"`
	Information  Information `json:"information"`
}

//Information holds some meta information for Ticker
type Information struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
}

//NewTicker creates new Ticker
func NewTicker() Ticker {
	return Ticker{
		CreationDate: time.Now(),
	}
}
