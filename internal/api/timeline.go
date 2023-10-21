package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/pagination"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/storage"
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
