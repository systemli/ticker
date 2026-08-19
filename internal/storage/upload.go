package storage

import (
	"fmt"
	"time"

	uuid2 "github.com/google/uuid"
)

// MediaPath is the path media is served on, as seen by a browser talking to the
// admin or frontend.
const MediaPath = "/api/media"

// uploadExtensions maps the content types the API accepts to the extension a file
// is stored and served under. The extension decides how a browser interprets the
// response and media shares an origin with the interfaces, so it is derived from
// the sniffed content type instead of the client-supplied filename.
var uploadExtensions = map[string]string{
	"image/gif":  "gif",
	"image/jpeg": "jpg",
	"image/png":  "png",
}

// ExtensionForContentType returns the extension uploads of this content type are
// stored under, and whether the type is accepted at all.
func ExtensionForContentType(contentType string) (string, bool) {
	ext, ok := uploadExtensions[contentType]
	return ext, ok
}

type Upload struct {
	ID          int `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UUID        string `gorm:"index;unique"`
	TickerID    int    `gorm:"index"`
	Path        string
	Extension   string
	ContentType string
}

func NewUpload(contentType string, tickerID int) Upload {
	now := time.Now()
	uuid := uuid2.New()
	ext, _ := ExtensionForContentType(contentType)
	// First version we use a date based directory structure
	path := fmt.Sprintf("%d/%d", now.Year(), now.Month())

	return Upload{
		Path:        path,
		UUID:        uuid.String(),
		TickerID:    tickerID,
		Extension:   ext,
		ContentType: contentType,
	}
}

func (u *Upload) FileName() string {
	return fmt.Sprintf("%s.%s", u.UUID, u.Extension)
}

func (u *Upload) FullPath(uploadPath string) string {
	return fmt.Sprintf("%s/%s/%s", uploadPath, u.Path, u.FileName())
}

func (u *Upload) URL() string {
	return MediaURL(u.FileName())
}

// MediaURL returns the public, host-relative URL of an uploaded file. The API
// serves media below /v1, which the admin and frontend reach through their own
// /api path, so no absolute base URL is needed.
func MediaURL(name string) string {
	return fmt.Sprintf("%s/%s", MediaPath, name)
}
