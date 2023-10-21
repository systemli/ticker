package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/gorm"
)

var log = logrus.WithField("package", "storage")

type Storage interface {
	CountUser() (int, error)
	FindUsers() ([]User, error)
	FindUserByID(id int) (User, error)
	FindUsersByIDs(ids []int) ([]User, error)
	FindUserByEmail(email string) (User, error)
	FindUsersByTicker(ticker Ticker) ([]User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error
	DeleteTickerUsers(ticker *Ticker) error
	DeleteTickerUser(ticker *Ticker, user *User) error
	AddTickerUser(ticker *Ticker, user *User) error
	FindTickers(opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error)
	FindTickersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error)
	FindTickerByDomain(domain string, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error)
	FindTickerByID(id int, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker Ticker) error
	SaveUpload(upload *Upload) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker Ticker) error
	FindMessage(tickerID, messageID int, opts ...func(*gorm.DB) *gorm.DB) (Message, error)
	FindMessagesByTicker(ticker Ticker, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error)
	FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error)
	SaveMessage(message *Message) error
	DeleteMessage(message Message) error
	DeleteMessages(ticker Ticker) error
	DeleteAttachmentsByMessage(message Message) error
	GetInactiveSettings() InactiveSettings
	GetRefreshIntervalSettings() RefreshIntervalSettings
	SaveInactiveSettings(inactiveSettings InactiveSettings) error
	SaveRefreshIntervalSettings(refreshInterval RefreshIntervalSettings) error
	UploadPath() string
}
