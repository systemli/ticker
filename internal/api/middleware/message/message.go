package message

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func PrefetchMessage(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticker, _ := helper.Ticker(c)

		messageID, err := strconv.Atoi(c.Param("messageID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
			return
		}

		message, err := s.FindMessage(ticker.ID, messageID, storage.WithAttachments())
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.MessageNotFound))
			return
		}

		c.Set("message", message)
	}
}
