package storage

// UploadStore covers Upload CRUD and the filesystem path where uploads live.
type UploadStore interface {
	SaveUpload(upload *Upload) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker *Ticker) error

	UploadPath() string
}
