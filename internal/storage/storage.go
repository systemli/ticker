package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/gorm"
)

var log = logrus.WithField("package", "storage")

type Storage interface {
	FindUsers(filter UserFilter, opts ...func(*gorm.DB) *gorm.DB) ([]User, error)
	FindUserByID(id int, opts ...func(*gorm.DB) *gorm.DB) (User, error)
	FindUsersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]User, error)
	FindUserByEmail(email string, opts ...func(*gorm.DB) *gorm.DB) (User, error)
	FindUsersByTicker(ticker Ticker, opts ...func(*gorm.DB) *gorm.DB) ([]User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error
	DeleteTickerUsers(ticker *Ticker) error
	DeleteTickerUser(ticker *Ticker, user *User) error
	AddTickerUser(ticker *Ticker, user *User) error
	FindTickersByUser(user User, filter TickerFilter, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error)
	FindTickerByUserAndID(user User, id int, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error)
	FindTickersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error)
	FindTickerByOrigin(origin string, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error)
	FindTickerByID(id int, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker *Ticker) error
	SaveTickerWebsites(ticker *Ticker, websites []TickerWebsite) error
	DeleteTickerWebsites(ticker *Ticker) error
	ResetTicker(ticker *Ticker) error
	DeleteIntegrations(ticker *Ticker) error
	DeleteMastodon(ticker *Ticker) error
	DeleteTelegram(ticker *Ticker) error
	DeleteBluesky(ticker *Ticker) error
	DeleteSignalGroup(ticker *Ticker) error
	SaveUpload(upload *Upload) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker *Ticker) error
	FindMessage(tickerID, messageID int, opts ...func(*gorm.DB) *gorm.DB) (Message, error)
	FindMessagesByTicker(ticker Ticker, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error)
	FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error)
	SaveMessage(message *Message) error
	DeleteMessage(message Message) error
	DeleteMessages(ticker *Ticker) error
	GetInactiveSettings() InactiveSettings
	GetRefreshIntervalSettings() RefreshIntervalSettings
	SaveInactiveSettings(inactiveSettings InactiveSettings) error
	SaveRefreshIntervalSettings(refreshInterval RefreshIntervalSettings) error
	UploadPath() string
}
