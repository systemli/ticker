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
	FindUserByEmail(email string) (User, error)
	FindUsersByTicker(ticker Ticker) ([]User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error
	AddUsersToTicker(ticker Ticker, ids []int) error
	RemoveTickerFromUser(ticker Ticker, user User) error
	FindTickers() ([]Ticker, error)
	FindTickersByIDs(ids []int) ([]Ticker, error)
	FindTickerByDomain(domain string) (Ticker, error)
	FindTickerByID(id int) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker Ticker) error
	FindUploadsByMessage(message Message) []Upload
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
	FindSetting(name string) (Setting, error)
	GetInactiveSetting() Setting
	GetRefreshIntervalSetting() Setting
	GetRefreshIntervalSettingValue() int
	SaveInactiveSetting(inactiveSettings InactiveSettings) error
	SaveRefreshInterval(refreshInterval float64) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	UploadPath() string
}
