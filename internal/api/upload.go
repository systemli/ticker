package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

var allowedContentTypes = []string{"image/jpeg", "image/gif", "image/png"}

func PostUpload(c *gin.Context) {
	me, err := Me(c)
	if checkError(c, err, http.StatusBadRequest, ErrorCodeDefault, ErrorUserNotFound) {
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
	if checkError(c, err, http.StatusBadRequest, ErrorCodeDefault, "can't convert ticker id to int") {
		return
	}
	ticker, err := GetTicker(tickerID)
	if checkError(c, err, http.StatusBadRequest, ErrorCodeDefault, ErrorTickerNotFound) {
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
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, ErrorTooMuchFiles))
		return
	}
	var uploads []*Upload
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if checkError(c, err, http.StatusBadRequest, ErrorCodeDefault, "can't open file in upload") {
			return
		}

		contentType := util.DetectContentType(file)
		if !util.ContainsString(allowedContentTypes, contentType) {
			log.Error(fmt.Sprintf("%s is not allowed to uploaded", contentType))
			c.JSON(http.StatusBadRequest, NewJSONErrorResponse(ErrorCodeDefault, "failed to upload"))
			return
		}

		u := NewUpload(fileHeader.Filename, contentType, ticker.ID)
		err = DB.Save(u)
		if checkError(c, err, http.StatusInternalServerError, ErrorCodeDefault, "can't save upload") {
			return
		}

		err = preparePath(u.FullPath())
		if checkError(c, err, http.StatusInternalServerError, ErrorCodeDefault, "can't prepare upload path") {
			return
		}

		if u.ContentType == "image/gif" {
			err = c.SaveUploadedFile(fileHeader, u.FullPath())
			if checkError(c, err, http.StatusInternalServerError, ErrorCodeDefault, "can't save gif") {
				return
			}
		} else {
			nFile, _ := fileHeader.Open()
			image, err := util.ResizeImage(nFile, 1280)
			if checkError(c, err, http.StatusInternalServerError, ErrorCodeDefault, "can't resize file") {
				return
			}

			err = util.SaveImage(image, u.FullPath())
			if checkError(c, err, http.StatusInternalServerError, ErrorCodeDefault, "can't save uploaded file") {
				return
			}
		}

		uploads = append(uploads, u)
	}

	c.JSON(http.StatusOK, NewJSONSuccessResponse("uploads", NewUploadsResponse(uploads)))
}

func checkError(c *gin.Context, err error, httpStatus, errorCode int, message string) bool {
	if err != nil {
		log.WithError(err).Error(message)
		c.JSON(httpStatus, NewJSONErrorResponse(errorCode, "failed to upload"))
		return true
	}

	return false
}

func preparePath(path string) error {
	fs := Config.FileBackend
	return fs.MkdirAll(filepath.Dir(path), 0750)
}
