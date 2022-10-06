package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/pagination"
	"github.com/systemli/ticker/internal/api/renderer"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetFeed(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	pagination := pagination.NewPagination(c)
	messages, err := h.storage.FindMessagesByTickerAndPagination(ticker, *pagination)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.MessageFetchError))
		return
	}

	format := renderer.FormatFromString(c.Query("format"))
	feed := buildFeed(ticker, messages)

	c.Render(http.StatusOK, renderer.Feed{Data: feed, Format: format})
}

func buildFeed(ticker storage.Ticker, messages []storage.Message) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       ticker.Title,
		Description: ticker.Description,
		Author: &feeds.Author{
			Name:  ticker.Information.Author,
			Email: ticker.Information.Email,
		},
		Link: &feeds.Link{
			Href: ticker.Information.URL,
		},
		Created: time.Now(),
	}

	items := make([]*feeds.Item, 0)
	for _, message := range messages {
		item := &feeds.Item{
			Id:          strconv.Itoa(message.ID),
			Created:     message.CreationDate,
			Description: message.Text,
			Title:       message.Text,
			Link:        &feeds.Link{},
		}
		items = append(items, item)
	}
	feed.Items = items

	return feed
}
