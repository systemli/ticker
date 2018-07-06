package model

import (
	"time"
)

//Message represents a single message
type Message struct {
	ID           int       `storm:"id,increment"`
	CreationDate time.Time `storm:"index"`
	Ticker       int       `storm:"index"`
	Text         string
	Tweet        Tweet
	//TODO: Geolocation, Facebook-ID
}

//
type Tweet struct {
	ID       string
	UserName string
}

type MessageResponse struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	Text         string    `json:"text"`
	Ticker       int       `json:"ticker"`
	TweetID      string    `json:"tweet_id"`
	TweetUser    string    `json:"tweet_user"`
}

//NewMessage creates new Message
func NewMessage() *Message {
	return &Message{
		CreationDate: time.Now(),
	}
}

//
func NewMessageResponse(message Message) *MessageResponse {
	return &MessageResponse{
		ID:           message.ID,
		CreationDate: message.CreationDate,
		Text:         message.Text,
		Ticker:       message.Ticker,
		TweetID:      message.Tweet.ID,
		TweetUser:    message.Tweet.UserName,
	}
}

//
func NewMessagesResponse(messages []Message) []*MessageResponse {
	var mr []*MessageResponse

	for _, message := range messages {
		mr = append(mr, NewMessageResponse(message))
	}

	return mr
}
