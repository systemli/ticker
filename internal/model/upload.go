package model

import (
	"fmt"
	"path/filepath"
	"time"

	uuid2 "github.com/google/uuid"
)

//Upload represents the structure of an Upload configuration
type Upload struct {
	ID           int       `storm:"id,increment"`
	UUID         string    `storm:"index,unique"`
	CreationDate time.Time `storm:"index"`
	TickerID     int       `storm:"index"`
	Path         string
	Extension    string
}

//UploadResponse represents the Upload for API responses.
type UploadResponse struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"`
	CreationDate time.Time `json:"creation_date"`
	URL          string    `json:"url"`
}

//NewUpload creates new Upload.
func NewUpload(filename string, tickerID int) *Upload {
	now := time.Now()
	uuid := uuid2.New()
	ext := filepath.Ext(filename)[1:]
	// First version we use a date based directory structure
	path := fmt.Sprintf("%d/%d", now.Year(), now.Month())

	return &Upload{
		CreationDate: now,
		Path:         path,
		UUID:         uuid.String(),
		TickerID:     tickerID,
		Extension:    ext,
	}
}

//FileName returns the name with file extension.
func (u *Upload) FileName() string {
	return fmt.Sprintf("%s.%s", u.UUID, u.Extension)
}

//URL returns the public url for the upload.
func (u *Upload) URL() string {
	return MediaURL(u.FileName())
}

//NewUploadResponse returns a API friendly representation for a Upload.
func NewUploadResponse(upload *Upload) *UploadResponse {
	return &UploadResponse{
		ID:           upload.ID,
		UUID:         upload.UUID,
		CreationDate: upload.CreationDate,
		URL:          upload.URL(),
	}
}

//NewTickersResponse prepares a map of []TickerResponse.
func NewUploadsResponse(uploads []*Upload) []*UploadResponse {
	var ur []*UploadResponse

	for _, upload := range uploads {
		ur = append(ur, NewUploadResponse(upload))
	}

	return ur
}

func MediaURL(name string) string {
	return fmt.Sprintf("%s/media/%s", Config.UploadURL, name)
}
