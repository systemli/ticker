package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

func PostUpload(c *gin.Context) {
	me, err := Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, NewJSONErrorResponse(ErrorCodeDefault, ErrorUserNotFound))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}

	if len(form.Value["ticker"]) != 1 {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorTickerIdentifierMissing))
		return
	}

	tickerID, err := strconv.Atoi(form.Value["ticker"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
		return
	}
	ticker, err := GetTicker(tickerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorTickerNotFound))
		return
	}

	if !me.IsSuperAdmin {
		if !contains(me.Tickers, tickerID) {
			c.JSON(http.StatusForbidden, NewJSONErrorResponse(ErrorCodeInsufficientPermissions, ErrorInsufficientPermissions))
			return
		}
	}

	files := form.File["files"]
	if len(files) < 1 {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorFilesIdentifierMissing))
		return
	}
	var uploads []*Upload
	for _, file := range files {
		u := NewUpload(file.Filename, ticker.ID)
		err = DB.Save(u)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		path := fmt.Sprintf("%s/%s", Config.UploadPath, u.Path)

		err := preparePath(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		dst := fmt.Sprintf("%s/%s", path, u.FileName())

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		uploads = append(uploads, u)
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("uploads", NewUploadsResponse(uploads)))
	return
}

func preparePath(path string) error {
	fs := Config.FileBackend
	return fs.MkdirAll(path, 0750)
}
