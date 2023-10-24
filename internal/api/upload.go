package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/systemli/ticker/internal/api/helper"
	"github.com/systemli/ticker/internal/api/response"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"

	"github.com/systemli/ticker/internal/util"
)

var allowedContentTypes = []string{"image/jpeg", "image/gif", "image/png"}

func (h *handler) PostUpload(c *gin.Context) {
	me, err := helper.Me(c)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(response.CodeDefault, response.UserNotFound))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
		return
	}

	if len(form.Value["ticker"]) != 1 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	tickerID, err := strconv.Atoi(form.Value["ticker"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerIdentifierMissing))
		return
	}

	ticker, err := h.storage.FindTickerByUserAndID(me, tickerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TickerNotFound))
		return
	}

	files := form.File["files"]
	if len(files) < 1 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FilesIdentifierMissing))
		return
	}
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.TooMuchFiles))
		return
	}
	uploads := make([]storage.Upload, 0)
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
			return
		}

		contentType := util.DetectContentType(file)
		if !util.ContainsString(allowedContentTypes, contentType) {
			log.Error(fmt.Sprintf("%s is not allowed to uploaded", contentType))
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, "failed to upload"))
			return
		}

		u := storage.NewUpload(fileHeader.Filename, contentType, ticker.ID)
		err = h.storage.SaveUpload(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
			return
		}

		err = preparePath(u, h.config)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse(response.CodeDefault, response.FormError))
			return
		}

		if u.ContentType == "image/gif" {
			err = c.SaveUploadedFile(fileHeader, u.FullPath(h.config.UploadPath))
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.FormError))
				return
			}
		} else {
			nFile, _ := fileHeader.Open()
			image, err := util.ResizeImage(nFile, 1280)
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.FormError))
				return
			}

			err = util.SaveImage(image, u.FullPath(h.config.UploadPath))
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.CodeDefault, response.FormError))
				return
			}
		}

		uploads = append(uploads, u)
	}

	c.JSON(http.StatusOK, response.SuccessResponse(map[string]interface{}{"uploads": response.UploadsResponse(uploads, h.config)}))
}

func preparePath(upload storage.Upload, config config.Config) error {
	path := upload.FullPath(config.UploadPath)
	fs := config.FileBackend
	return fs.MkdirAll(filepath.Dir(path), 0750)
}
