package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetMedia(c *gin.Context) {
	parts := strings.Split(c.Param("fileName"), ".")
	upload, err := h.storage.FindUploadByUUID(parts[0])
	if err != nil {
		c.String(http.StatusNotFound, "%s", err.Error())
		return
	}

	expireTime := time.Now().AddDate(0, 1, 0)
	cacheControl := fmt.Sprintf("public, max-age=%d", expireTime.Unix())
	expires := expireTime.Format(http.TimeFormat)

	c.Header("Cache-Control", cacheControl)
	c.Header("Expires", expires)
	c.File(upload.FullPath(h.storage.UploadPath()))
}
