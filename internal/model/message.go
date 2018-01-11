package model

import "time"

//Message represents a single message
type Message struct {
	ID           int       `json:"id" storm:"id,increment"`
	CreationDate time.Time `json:"creation_date" storm:"index"`
	Text         string    `json:"text" binding:"required"`
	Ticker       int       `json:"ticker" binding:"required" storm:"index"`
	//TODO: Geolocation, Tweet/Facebook-ID
}

//NewMessage creates new Message
func NewMessage() Message {
	return Message{
		CreationDate: time.Now(),
	}
}
