package model

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	geojson "github.com/paulmach/go.geojson"
)

//Message represents a single message
type Message struct {
	ID             int       `storm:"id,increment"`
	CreationDate   time.Time `storm:"index"`
	Ticker         int       `storm:"index"`
	Text           string
	Attachments    []Attachment
	GeoInformation geojson.FeatureCollection
	Tweet          Tweet
	Telegram       TelegramMeta
	//TODO: Facebook-ID
}

//
type Tweet struct {
	ID       string
	UserName string
}

type TelegramMeta struct {
	Messages []tgbotapi.Message
}

type Attachment struct {
	UUID        string
	Extension   string
	ContentType string
}

type MessageResponse struct {
	ID               int                          `json:"id"`
	CreationDate     time.Time                    `json:"creation_date"`
	Text             string                       `json:"text"`
	Ticker           int                          `json:"ticker"`
	TweetID          string                       `json:"tweet_id"`
	TweetUser        string                       `json:"tweet_user"`
	TelegramMessages []tgbotapi.Message           `json:"telegram_messages"`
	GeoInformation   string                       `json:"geo_information"`
	Attachments      []*MessageAttachmentResponse `json:"attachments"`
}

type MessageAttachmentResponse struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

//NewMessage creates new Message
func NewMessage() *Message {
	return &Message{
		CreationDate: time.Now(),
	}
}

//
func NewMessageResponse(message Message) *MessageResponse {
	m, _ := message.GeoInformation.MarshalJSON()
	var attachments []*MessageAttachmentResponse

	for _, attachment := range message.Attachments {
		name := fmt.Sprintf("%s.%s", attachment.UUID, attachment.Extension)
		attachments = append(attachments, &MessageAttachmentResponse{URL: MediaURL(name), ContentType: attachment.ContentType})
	}

	return &MessageResponse{
		ID:               message.ID,
		CreationDate:     message.CreationDate,
		Text:             message.Text,
		Ticker:           message.Ticker,
		TweetID:          message.Tweet.ID,
		TweetUser:        message.Tweet.UserName,
		TelegramMessages: message.Telegram.Messages,
		GeoInformation:   string(m),
		Attachments:      attachments,
	}
}

//
func NewMessagesResponse(messages []Message) []*MessageResponse {
	mr := make([]*MessageResponse, 0)
	for _, message := range messages {
		mr = append(mr, NewMessageResponse(message))
	}

	return mr
}
