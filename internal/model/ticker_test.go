package model_test

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestTwitter_Connected(t *testing.T) {
	ticker := model.NewTicker()

	assert.False(t, ticker.Twitter.Connected())

	ticker.Twitter.Secret = "secret"
	ticker.Twitter.Token = "token"

	assert.True(t, ticker.Twitter.Connected())
}

func TestTicker_Reset(t *testing.T) {
	ticker := model.NewTicker()

	ticker.Active = true
	ticker.Description = "description"
	ticker.PrependTime = true
	ticker.Hashtags = []string{"hashtag"}
	ticker.Information = model.Information{
		Author: "author",
	}
	ticker.Twitter.Secret = "secret"
	ticker.Twitter.Token = "token"
	ticker.Twitter.Active = true
	ticker.Twitter.User = twitter.User{}
	ticker.Location = model.Location{}

	assert.True(t, ticker.Active)

	ticker.Reset()

	assert.False(t, ticker.Active)
	assert.Equal(t, "", ticker.Description)
	assert.False(t, ticker.PrependTime)
	assert.Equal(t, []string{}, ticker.Hashtags)
	assert.Equal(t, model.Information{}, ticker.Information)
	assert.Equal(t, "", ticker.Twitter.Secret)
	assert.Equal(t, "", ticker.Twitter.Token)
	assert.False(t, ticker.Twitter.Active)
	assert.Equal(t, twitter.User{}, ticker.Twitter.User)
	assert.Equal(t, model.Location{}, ticker.Location)
}

func TestNewTickerResponse(t *testing.T) {
	ticker := model.NewTicker()
	r := model.NewTickerResponse(ticker)

	assert.Equal(t, 0, r.ID)
}

func TestNewTickersResponse(t *testing.T) {
	ticker := model.NewTicker()
	r := model.NewTickersResponse([]*model.Ticker{ticker})

	assert.Equal(t, 1, len(r))
}
