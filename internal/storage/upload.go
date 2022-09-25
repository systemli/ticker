package storage

import (
	"fmt"
	"path/filepath"
	"time"

	uuid2 "github.com/google/uuid"
)

type Upload struct {
	ID           int       `storm:"id,increment"`
	UUID         string    `storm:"index,unique"`
	CreationDate time.Time `storm:"index"`
	TickerID     int       `storm:"index"`
	Path         string
	Extension    string
	ContentType  string
}

func NewUpload(filename, contentType string, tickerID int) Upload {
	now := time.Now()
	uuid := uuid2.New()
	ext := filepath.Ext(filename)[1:]
	// First version we use a date based directory structure
	path := fmt.Sprintf("%d/%d", now.Year(), now.Month())

	return Upload{
		CreationDate: now,
		Path:         path,
		UUID:         uuid.String(),
		TickerID:     tickerID,
		Extension:    ext,
		ContentType:  contentType,
	}
}

func (u *Upload) FileName() string {
	return fmt.Sprintf("%s.%s", u.UUID, u.Extension)
}

func (u *Upload) FullPath(uploadPath string) string {
	return fmt.Sprintf("%s/%s/%s", uploadPath, u.Path, u.FileName())
}

func (u *Upload) URL(uploadPath string) string {
	return MediaURL(u.FileName(), uploadPath)
}

func MediaURL(name string, uploadPath string) string {
	return fmt.Sprintf("%s/media/%s", uploadPath, name)
}
