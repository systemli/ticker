package response

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type InitTickerResponseTestSuite struct {
	suite.Suite
}

func (s *InitTickerResponseTestSuite) TestInitTickerResponse() {
	ticker := storage.Ticker{
		ID:          1,
		CreatedAt:   time.Now(),
		Domain:      "example.com",
		Title:       "Example",
		Description: "Example",
		Information: storage.TickerInformation{
			Author:    "Example",
			URL:       "https://example.com",
			Email:     "contact@example.com",
			Twitter:   "example",
			Facebook:  "example",
			Instagram: "example",
			Threads:   "example",
			Telegram:  "example",
			Mastodon:  "example",
			Bluesky:   "example",
		},
	}

	response := InitTickerResponse(ticker)

	s.Equal(ticker.ID, response.ID)
	s.Equal(ticker.CreatedAt, response.CreatedAt)
	s.Equal(ticker.Title, response.Title)
	s.Equal(ticker.Description, response.Description)
	s.Equal(ticker.Information.Author, response.Information.Author)
	s.Equal(ticker.Information.URL, response.Information.URL)
	s.Equal(ticker.Information.Email, response.Information.Email)
	s.Equal(ticker.Information.Twitter, response.Information.Twitter)
	s.Equal(ticker.Information.Facebook, response.Information.Facebook)
	s.Equal(ticker.Information.Instagram, response.Information.Instagram)
	s.Equal(ticker.Information.Threads, response.Information.Threads)
	s.Equal(ticker.Information.Telegram, response.Information.Telegram)
	s.Equal(ticker.Information.Mastodon, response.Information.Mastodon)
	s.Equal(ticker.Information.Bluesky, response.Information.Bluesky)
}

func TestInitTickerResponseTestSuite(t *testing.T) {
	suite.Run(t, new(InitTickerResponseTestSuite))
}
