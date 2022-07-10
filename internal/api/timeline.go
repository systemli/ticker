package api

import (
	"fmt"
	"net/http"
	"time"

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
	Messages []message `json:"messages"`
}

type message struct {
	ID             int          `json:"id"`
	CreationDate   time.Time    `json:"creation_date"`
	Text           string       `json:"text"`
	GeoInformation string       `json:"geo_information"`
	Attachments    []attachment `json:"attachments"`
}

type attachment struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

func newTimelineResponse(msgs []Message) *timelineResponse {
	messages := make([]message, 0)

	for _, msg := range msgs {
		var attachments []attachment

		geoInformation, _ := msg.GeoInformation.MarshalJSON()

		for _, a := range msg.Attachments {
			name := fmt.Sprintf("%s.%s", a.UUID, a.Extension)
			attachments = append(attachments, attachment{URL: MediaURL(name), ContentType: a.ContentType})
		}

		messages = append(messages, message{
			ID:             msg.ID,
			CreationDate:   msg.CreationDate,
			Text:           msg.Text,
			GeoInformation: string(geoInformation),
			Attachments:    attachments,
		})
	}

	return &timelineResponse{
		Data: timelineData{
			Messages: messages,
		},
		Status: ResponseSuccess,
	}
}

// GetTimelineHandler returns the public timeline for a ticker.
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

	c.JSON(http.StatusOK, newTimelineResponse(messages))
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
