package response

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type Upload struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"createdAt"`
	URL         string    `json:"url"`
	ContentType string    `json:"contentType"`
}

func UploadResponse(upload storage.Upload) Upload {
	return Upload{
		ID:          upload.ID,
		UUID:        upload.UUID,
		CreatedAt:   upload.CreatedAt,
		URL:         upload.URL(),
		ContentType: upload.ContentType,
	}
}

func UploadsResponse(uploads []storage.Upload) []Upload {
	ur := make([]Upload, 0)
	for _, upload := range uploads {
		ur = append(ur, UploadResponse(upload))
	}

	return ur
}
