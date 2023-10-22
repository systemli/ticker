package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/pagination"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func (h *handler) GetTimeline(c *gin.Context) {
	ticker, err := helper.Ticker(c)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	messages := make([]storage.Message, 0)
	if ticker.Active {
		pagination := pagination.NewPagination(c)
		messages, err = h.storage.FindMessagesByTickerAndPagination(ticker, *pagination, storage.WithAttachments())
		if err != nil {
			c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.MessageFetchError))
			return
		}
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"messages": response.TimelineResponse(messages, h.config)}))
}
