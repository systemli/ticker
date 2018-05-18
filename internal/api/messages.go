package api

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)

//GetMessages returns all Messages with paging
func GetMessages(c *gin.Context) {
	var ticker Ticker

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, "ticker not found"))
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

		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("messages", messages))
}

//GetMessage returns a Message for the given id
func GetMessage(c *gin.Context) {
	var ticker Ticker

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, "ticker not found"))
		return
	}

	var message Message

	messageID, err := strconv.Atoi(c.Param("messageID"))

	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", message))
}

//PostMessage creates and returns a new Message
func PostMessage(c *gin.Context) {
	var ticker Ticker

	message := NewMessage()
	err := c.Bind(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	message.Ticker = tickerID

	err = DB.Save(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", message))
}

//DeleteTicker deletes a existing Ticker
func DeleteMessage(c *gin.Context) {
	var ticker Ticker

	tickerID, err := strconv.Atoi(c.Param("tickerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", tickerID, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	messageID, err := strconv.Atoi(c.Param("messageID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	var message Message

	err = DB.One("ID", messageID, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.DeleteStruct(&message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   nil,
		"status": ResponseSuccess,
		"error":  nil,
	})
}
