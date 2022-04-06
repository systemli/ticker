package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"

	"github.com/systemli/ticker/internal/bridge"
	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

//GetMessagesHandler returns all Messages with paging
func GetMessagesHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, ErrorTickerNotFound))
		return
	}

	var messages []Message
	//TODO: Pagination
	err = DB.Find("Ticker", tickerID, &messages, storm.Reverse())
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusOK, NewJSONSuccessResponse("messages", []string{}))
			return
		}

		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("messages", NewMessagesResponse(messages)))
}

//GetMessageHandler returns a Message for the given id
func GetMessageHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, ErrorTickerNotFound))
		return
	}

	var message Message
	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}
	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", NewMessageResponse(message)))
}

//PostMessageHandler creates and returns a new Message
func PostMessageHandler(c *gin.Context) {
	var body struct {
		Text           string                    `json:"text" binding:"required"`
		GeoInformation geojson.FeatureCollection `json:"geo_information"`
		Attachments    []int                     `json:"attachments"`
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	ticker := NewTicker()
	err = DB.One("ID", tickerID, ticker)
	if err != nil {
		log.WithError(err).Error("failed to find the ticker")
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, "ticker not found"))
		return
	}

	var uploads []*Upload
	if len(body.Attachments) > 0 {
		err := DB.Select(q.In("ID", body.Attachments)).Find(&uploads)
		if err != nil {
			c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
			return
		}
	}

	message := NewMessage()
	message.Text = body.Text
	message.Ticker = tickerID
	message.GeoInformation = body.GeoInformation

	if len(uploads) > 0 {
		var attachments []Attachment
		for _, upload := range uploads {
			attachments = append(attachments, Attachment{Extension: upload.Extension, UUID: upload.UUID, ContentType: upload.ContentType})
		}

		message.Attachments = attachments
	}

	err = bridge.SendTweet(ticker, message)
	if err != nil {
		log.WithError(err).WithField("ticker", ticker.ID).WithField("message", message.ID).Error("sending message to twitter failed")
	}

	err = DB.Save(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", NewMessageResponse(*message)))
}

//DeleteTickerHandler deletes a existing Ticker
func DeleteMessageHandler(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	var ticker Ticker
	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeNotFound, err.Error()))
		return
	}

	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	var message Message

	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	err = DeleteMessage(&ticker, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}
	err = bridge.DeleteTweet(&ticker, &message)
	if err != nil {
		log.WithField("error", err).WithField("message", message).Error("failed to delete tweet")
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}
