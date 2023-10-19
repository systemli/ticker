package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/pagination"
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
	FindTickers() ([]Ticker, error)
	FindTickersByIDs(ids []int) ([]Ticker, error)
	FindTickerByDomain(domain string) (Ticker, error)
	FindTickerByID(id int) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker Ticker) error
	SaveUpload(upload *Upload) error
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker Ticker) error
	FindMessage(tickerID, messageID int) (Message, error)
	FindMessagesByTicker(ticker Ticker) ([]Message, error)
	FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination) ([]Message, error)
	SaveMessage(message *Message) error
	DeleteMessage(message Message) error
	DeleteMessages(ticker Ticker) error
	GetInactiveSettings() InactiveSettings
	GetRefreshIntervalSettings() RefreshIntervalSettings
	SaveInactiveSettings(inactiveSettings InactiveSettings) error
	SaveRefreshIntervalSettings(refreshInterval RefreshIntervalSettings) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	UploadPath() string
}
