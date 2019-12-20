package api

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

func GetMedia(c *gin.Context) {
	var upload Upload

	parts := strings.Split(c.Param("fileName"), ".")
	err := DB.One("UUID", parts[0], &upload)
	if err != nil {
		c.String(http.StatusNotFound, "%s", err.Error())
		return
	}

	file, err := Config.FileBackend.Open(upload.FullPath())
	if err != nil {
		c.String(http.StatusInternalServerError, "serve error: %s", err.Error())
		return
	}

	stat, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "serve error: %s", err.Error())
		return
	}

	contentType := util.DetectContentType(file)
	expireTime := time.Now().AddDate(0, 1, 0)
	cacheControl := fmt.Sprintf("public, max-age=%d", expireTime.Unix())
	expires := expireTime.Format(http.TimeFormat)
	reader := bufio.NewReader(file)

	c.Header("Cache-Control", cacheControl)
	c.Header("Expires", expires)
	c.DataFromReader(http.StatusOK, stat.Size(), contentType, reader, map[string]string{})
}
