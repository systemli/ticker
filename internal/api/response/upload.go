package response

import (
	"time"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type Upload struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"createdAt"`
	URL         string    `json:"url"`
	ContentType string    `json:"contentType"`
}

func UploadResponse(upload storage.Upload, config config.Config) Upload {
	return Upload{
		ID:          upload.ID,
		UUID:        upload.UUID,
		CreatedAt:   upload.CreatedAt,
		URL:         upload.URL(config.UploadURL),
		ContentType: upload.ContentType,
	}
}

func UploadsResponse(uploads []storage.Upload, config config.Config) []Upload {
	ur := make([]Upload, 0)
	for _, upload := range uploads {
		ur = append(ur, UploadResponse(upload, config))
	}

	return ur
}
