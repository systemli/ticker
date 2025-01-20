package storage

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ticker = NewTicker()

func TestTickerMastodonConnect(t *testing.T) {
	assert.False(t, ticker.Mastodon.Connected())
}

func TestTickerTelegramConnected(t *testing.T) {
	assert.False(t, ticker.Telegram.Connected())

	ticker.Telegram.ChannelName = "ChannelName"

	assert.True(t, ticker.Telegram.Connected())
}

func TestTickerBlueskyConnected(t *testing.T) {
	assert.False(t, ticker.Bluesky.Connected())

	ticker.Bluesky.Handle = "Handle"
	ticker.Bluesky.AppKey = "AppKey"

	assert.True(t, ticker.Bluesky.Connected())
}

func TestTickerSignalGroupConnect(t *testing.T) {
	assert.False(t, ticker.SignalGroup.Connected())

	ticker.SignalGroup.GroupID = "GroupID"

	assert.True(t, ticker.SignalGroup.Connected())
}

func TestNewTickerFilter(t *testing.T) {
	filter := NewTickerFilter(nil)
	assert.Nil(t, filter.Active)
	assert.Nil(t, filter.Domain)
	assert.Nil(t, filter.Title)

	req := httptest.NewRequest("GET", "/", nil)
	filter = NewTickerFilter(req)
	assert.Nil(t, filter.Active)
	assert.Nil(t, filter.Domain)
	assert.Nil(t, filter.Title)

	req = httptest.NewRequest("GET", "/?active=true&domain=example.org&title=Title", nil)
	filter = NewTickerFilter(req)
	assert.True(t, *filter.Active)
	assert.Equal(t, "example.org", *filter.Domain)
	assert.Equal(t, "Title", *filter.Title)

	req = httptest.NewRequest("GET", "/?order_by=created_at&sort=asc", nil)
	filter = NewTickerFilter(req)
	assert.Equal(t, "created_at", filter.OrderBy)
	assert.Equal(t, "asc", filter.Sort)
}
