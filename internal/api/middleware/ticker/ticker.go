package ticker

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
)

func PrefetchTicker(s storage.TickerStore, opts ...storage.QueryOpt) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := helper.Me(c)
		tickerID, err := strconv.Atoi(c.Param("tickerID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
			return
		}

		ticker, err := s.FindTickerByUserAndID(user, tickerID, opts...)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.TickerNotFound))
			return
		}

		c.Set("ticker", ticker)
	}
}

func PrefetchTickerFromRequest(s storage.TickerStore, opts ...storage.QueryOpt) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, err := helper.GetOrigin(c)
		if err != nil {
			c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
			return
		}

		ticker, err := s.FindTickerByOrigin(origin, opts...)
		if err != nil {
			c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
			return
		}

		c.Set("ticker", ticker)
	}
}
