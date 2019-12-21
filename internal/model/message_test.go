package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestPrepareTweet(t *testing.T) {
	ticker := model.NewTicker()
	message := model.NewMessage()
	message.CreationDate, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	message.Text = "example"

	assert.Equal(t, "example", message.PrepareTweet(ticker))

	ticker.PrependTime = true

	assert.Equal(t, "22:08 example", message.PrepareTweet(ticker))
}

func TestNewMessageResponse(t *testing.T) {
	m := model.NewMessage()
	r := model.NewMessageResponse(*m)

	assert.Equal(t, 0, r.ID)
	assert.Equal(t, "", r.Text)
	assert.Equal(t, 0, r.Ticker)
	assert.Equal(t, "", r.TweetID)
	assert.Equal(t, "", r.TweetUser)
	assert.Equal(t, `{"type":"FeatureCollection","features":[]}`, r.GeoInformation)
}

func TestNewMessagesResponse(t *testing.T) {
	m := model.NewMessage()
	r := model.NewMessagesResponse([]model.Message{*m})

	assert.Equal(t, 1, len(r))
}
