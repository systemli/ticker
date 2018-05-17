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
	var messages []Message

	if len(c.Query("ticker")) == 0 {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, `Missing parameter "ticker"`))
		return
	}

	id, err := strconv.Atoi(c.Query("ticker"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	//TODO: Pagination
	err = DB.Find("Ticker", id, &messages, storm.Reverse())
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
	var message Message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", message))
}

//PostMessage creates and returns a new Message
func PostMessage(c *gin.Context) {
	message := NewMessage()
	err := c.Bind(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	//TODO: Find Ticker for ID
	if message.Ticker == 0 {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.Save(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", message))
}

//PutTicker updates and returns a existing Ticker
func PutMessage(c *gin.Context) {
	var message Message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = c.Bind(&message)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.Update(&message)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("message", message))
}

//DeleteTicker deletes a existing Ticker
func DeleteMessage(c *gin.Context) {
	var message Message
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &message)
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