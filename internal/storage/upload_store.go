package storage

import "gorm.io/gorm"

// UploadStore covers Upload CRUD and the filesystem path where uploads live.
type UploadStore interface {
	SaveUpload(upload *Upload) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker *Ticker) error

	UploadPath() string

	WithUploadTx(tx *gorm.DB) UploadStore
}

// WithUploadTx returns an UploadStore scoped to the given transaction.
func (s *SqlStorage) WithUploadTx(tx *gorm.DB) UploadStore {
	return &SqlStorage{DB: tx, uploadPath: s.uploadPath}
}
