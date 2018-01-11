package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/storage"
)


//GetTickers returns all Ticker with paging
func GetTickers(c *gin.Context) {
	var tickers []Ticker

	//TODO: Pagination
	err := DB.All(&tickers)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("tickers", tickers))
}

//GetTicker returns a Ticker for the given id
func GetTicker(c *gin.Context) {
	var ticker Ticker
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", ticker))
}

//PostTicker creates and returns a new Ticker
func PostTicker(c *gin.Context) {
	ticker := NewTicker()
	err := c.Bind(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.Save(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", ticker))
}

//PutTicker updates and returns a existing Ticker
func PutTicker(c *gin.Context) {
	var ticker Ticker
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = c.Bind(&ticker)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.Update(&ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("ticker", ticker))
}

//DeleteTicker deletes a existing Ticker
func DeleteTicker(c *gin.Context) {
	var ticker Ticker
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.One("ID", id, &ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorUnspecified, err.Error()))
		return
	}

	err = DB.DeleteStruct(&ticker)
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
