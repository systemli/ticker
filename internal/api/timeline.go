package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
	. "github.com/systemli/ticker/internal/util"
)

type timelineResponse struct {
	Data   timelineData `json:"data"`
	Status string       `json:"status"`
}

type timelineData struct {
	Messages []*MessageResponse `json:"messages"`
}

// GetTimelineHandler returns the public timeline for a ticker.
// @Summary      Fetch the messages for a ticker.
// @Description  Endpoint to retrieve the messages from a ticker. The endpoint has a pagination to fetch newer or older
// @Description  messages.
// @Tags         public
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit for fetched messages, default: 10"
// @Param        before  query     int  false  "ID of the message we look for older entries"
// @Param        after   query     int  false  "ID of the message we look for newer entries"
// @Success      200     {object}  timelineResponse
// @Success      400     {object}  errorResponse
// @Failure      500     {object}  interface{}
// @Router       /timeline [get]
func GetTimelineHandler(c *gin.Context) {
	domain, err := GetDomain(c)
	if err != nil {
		timelineErrorResponse(c, "Could not find a ticker.")
		return
	}

	ticker, err := FindTicker(domain)
	if err != nil {
		timelineErrorResponse(c, "Could not find a ticker.")
		return
	}

	pagination := NewPagination(c)
	messages, err := FindByTicker(ticker, pagination)
	if err != nil {
		timelineErrorResponse(c, "Could not load messages.")
		return
	}

	c.JSON(http.StatusOK, timelineResponse{
		Data:   timelineData{Messages: NewMessagesResponse(messages)},
		Status: ResponseSuccess,
	})
}

func timelineErrorResponse(c *gin.Context, m string) {
	c.JSON(http.StatusBadRequest, errorResponse{
		Data:   nil,
		Status: ResponseError,
		Error: errorData{
			Code:    ErrorCodeDefault,
			Message: m,
		},
	})
}
