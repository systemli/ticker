package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetMedia(c *gin.Context) {
	parts := strings.Split(c.Param("fileName"), ".")
	upload, err := h.storage.FindUploadByUUID(parts[0])
	if err != nil {
		c.String(http.StatusNotFound, "%s", err.Error())
		return
	}

	// Media is served on the same origin as the admin and frontend. The upload
	// handler only accepts JPEG, GIF and PNG, but be explicit about the type and
	// forbid sniffing anyway.
	c.Header("Content-Type", upload.ContentType)
	c.Header("X-Content-Type-Options", "nosniff")
	// Rows created before the extension was derived from the content type may
	// still carry an arbitrary one, so neutralise them explicitly.
	c.Header("Content-Security-Policy", "default-src 'none'; sandbox")
	c.Header("Content-Disposition", `inline; filename="`+upload.FileName()+`"`)
	// File names contain a UUID, so a response never becomes stale.
	c.Header("Cache-Control", "public, max-age=2592000, immutable")
	c.File(upload.FullPath(h.storage.UploadPath()))
}
