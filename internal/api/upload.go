package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

var allowedContentTypes = []string{"image/jpeg", "image/gif", "image/png"}

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
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		contentType := util.DetectContentType(file)
		if !util.ContainsString(allowedContentTypes, contentType) {
			c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, fmt.Sprintf("%s is not allowed to uploaded", contentType)))
			return
		}

		u := NewUpload(fileHeader.Filename, contentType, ticker.ID)
		err = DB.Save(u)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		err = preparePath(u.FullPath())
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		if err := c.SaveUploadedFile(fileHeader, u.FullPath()); err != nil {
			c.JSON(http.StatusInternalServerError, NewJSONErrorResponse(ErrorCodeDefault, err.Error()))
			return
		}

		uploads = append(uploads, u)
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("uploads", NewUploadsResponse(uploads)))
}

func preparePath(path string) error {
	fs := Config.FileBackend
	return fs.MkdirAll(filepath.Dir(path), 0750)
}
