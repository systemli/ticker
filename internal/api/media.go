package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/storage"
)

func GetMedia(c *gin.Context) {
	var upload Upload

	parts := strings.Split(c.Param("fileName"), ".")
	err := DB.One("UUID", parts[0], &upload)
	if err != nil {
		c.String(http.StatusNotFound, "%s", err.Error())
		return
	}

	expireTime := time.Now().AddDate(0, 1, 0)
	cacheControl := fmt.Sprintf("public, max-age=%d", expireTime.Unix())
	expires := expireTime.Format(http.TimeFormat)

	c.Header("Cache-Control", cacheControl)
	c.Header("Expires", expires)
	c.File(upload.FullPath())
}
