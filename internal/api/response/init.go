package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type InitTicker struct {
	ID          int                   `json:"id"`
	CreatedAt   time.Time             `json:"createdAt"`
	Domain      string                `json:"domain"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Information InitTickerInformation `json:"information"`
}

type InitTickerInformation struct {
	Author   string `json:"author"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
	Telegram string `json:"telegram"`
}

func InitTickerResponse(ticker storage.Ticker) InitTicker {
	return InitTicker{
		ID:          ticker.ID,
		CreatedAt:   ticker.CreatedAt,
		Domain:      ticker.Domain,
		Title:       ticker.Title,
		Description: ticker.Description,
		Information: InitTickerInformation{
			Author:   ticker.Information.Author,
			URL:      ticker.Information.URL,
			Email:    ticker.Information.Email,
			Twitter:  ticker.Information.Twitter,
			Facebook: ticker.Information.Facebook,
			Telegram: ticker.Information.Telegram,
		},
	}
}
