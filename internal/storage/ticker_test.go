package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ticker = NewTicker()

func TestTickerMastodonConnect(t *testing.T) {
	assert.False(t, ticker.Mastodon.Connected())
}

func TestTickerReset(t *testing.T) {
	ticker.Active = true
	ticker.Description = "Description"
	ticker.Information.Author = "Author"
	ticker.Information.Email = "Email"
	ticker.Information.Twitter = "Twitter"
	ticker.Telegram.Active = true
	ticker.Telegram.ChannelName = "ChannelName"
	ticker.Location.Lat = 1
	ticker.Location.Lon = 2

	ticker.Reset()

	assert.False(t, ticker.Active)
	assert.False(t, ticker.Telegram.Active)
	assert.Empty(t, ticker.Description)
	assert.Empty(t, ticker.Information.Author)
	assert.Empty(t, ticker.Information.Email)
	assert.Empty(t, ticker.Information.Twitter)
	assert.Empty(t, ticker.Telegram.ChannelName)
	assert.Empty(t, ticker.Location)
}
