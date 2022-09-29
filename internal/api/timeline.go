package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/pagination"
	"github.com/systemli/ticker/internal/api/response"
)

// GetTimeline returns the public timeline for a ticker.
// @Summary      Fetch the messages for a ticker.
// @Description  Endpoint to retrieve the messages from a ticker. The endpoint has a pagination to fetch newer or older
// @Description  messages. It is mandatory that the browser sends the origin as a header. This can be overwritten with
// @Description  a query parameter.
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        origin  query     string  false  "Origin from the ticker, e.g. demoticker.org"
// @Param        origin  header    string  false  "Origin from the ticker, e.g. http://demoticker.org"
// @Param        limit   query     int     false  "Limit for fetched messages, default: 10"
// @Param        before  query     int     false  "ID of the message we look for older entries"
// @Param        after   query     int     false  "ID of the message we look for newer entries"
// @Success      200     {object}  response.Response
// @Success      400     {object}  response.Response
// @Failure      500     {object}  response.Response
// @Router       /timeline [get]
func (h *handler) GetTimeline(c *gin.Context) {
	domain, err := helper.GetDomain(c)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	ticker, err := h.storage.FindTickerByDomain(domain)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	pagination := pagination.NewPagination(c)
	messages, err := h.storage.FindMessagesByTicker(ticker, *pagination)
	if err != nil {
		c.JSON(http.StatusOK, response.ErrorResponse(response.CodeDefault, response.MessageFetchError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"messages": response.MessagesResponse(messages, h.config)}))
}
