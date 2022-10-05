package ticker

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

func PrefetchTicker(s storage.TickerStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := helper.Me(c)
		tickerID, err := strconv.Atoi(c.Param("tickerID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
			return
		}

		if !user.IsSuperAdmin {
			if !util.Contains(user.Tickers, tickerID) {
				c.JSON(http.StatusForbidden, response.ErrorResponse(response.CodeInsufficientPermissions, response.InsufficientPermissions))
				return
			}
		}

		ticker, err := s.FindTickerByID(tickerID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeNotFound, response.TickerNotFound))
			return
		}

		c.Set("ticker", ticker)
	}
}

func PrefetchTickerFromRequest(s storage.TickerStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain, err := helper.GetDomain(c)
		if err != nil {
			c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
			return
		}

		ticker, err := s.FindTickerByDomain(domain)
		if err != nil {
			c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
			return
		}

		c.Set("ticker", ticker)
	}
}
