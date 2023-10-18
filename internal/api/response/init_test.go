package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/storage"
)

func TestInitTickerResponse(t *testing.T) {
	ticker := storage.Ticker{
		ID:          1,
		CreatedAt:   time.Now(),
		Domain:      "example.com",
		Title:       "Example",
		Description: "Example",
		Information: storage.TickerInformation{
			Author:   "Example",
			URL:      "https://example.com",
			Email:    "contact@example.com",
			Twitter:  "example",
			Facebook: "example",
			Telegram: "example",
		},
	}

	response := InitTickerResponse(ticker)

	assert.Equal(t, ticker.ID, response.ID)
	assert.Equal(t, ticker.CreatedAt, response.CreatedAt)
	assert.Equal(t, ticker.Domain, response.Domain)
	assert.Equal(t, ticker.Title, response.Title)
	assert.Equal(t, ticker.Description, response.Description)
	assert.Equal(t, ticker.Information.Author, response.Information.Author)
	assert.Equal(t, ticker.Information.URL, response.Information.URL)
	assert.Equal(t, ticker.Information.Email, response.Information.Email)
	assert.Equal(t, ticker.Information.Twitter, response.Information.Twitter)
	assert.Equal(t, ticker.Information.Facebook, response.Information.Facebook)
	assert.Equal(t, ticker.Information.Telegram, response.Information.Telegram)
}
