package storage

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pagination "github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSqlStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SqlStorage Suite")
}

var _ = Describe("SqlStorage", func() {
	db, err := gorm.Open(sqlite.Open("ticker.db"), &gorm.Config{})
	Expect(err).ToNot(HaveOccurred())

	var storage = NewSqlStorage(db, "/uploads")

	err = db.AutoMigrate(
		&Ticker{},
		&TickerInformation{},
		&TickerTelegram{},
		&TickerMastodon{},
		&User{},
		&Message{},
		&Upload{},
		&Attachment{},
		&Setting{},
	)
	Expect(err).ToNot(HaveOccurred())

	BeforeEach(func() {
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM messages")
		db.Exec("DELETE FROM attachments")
		db.Exec("DELETE FROM tickers")
		db.Exec("DELETE FROM settings")
		db.Exec("DELETE FROM uploads")
	})

	Describe("CountUser", func() {
		It("returns the number of users", func() {
			Expect(storage.CountUser()).To(Equal(0))

			err := db.Create(&User{}).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(storage.CountUser()).To(Equal(1))
		})
	})

	Describe("FindUsers", func() {
		It("returns all users", func() {
			users, err := storage.FindUsers()
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(BeEmpty())

			err = db.Create(&User{}).Error
			Expect(err).ToNot(HaveOccurred())

			users, err = storage.FindUsers()
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(HaveLen(1))
		})
	})

	Describe("FindUserByID", func() {
		It("returns the user with the given id", func() {
			user, err := storage.FindUserByID(1)
			Expect(err).To(HaveOccurred())
			Expect(user).To(BeZero())

			err = db.Create(&User{}).Error
			Expect(err).ToNot(HaveOccurred())

			user, err = storage.FindUserByID(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeZero())
		})
	})

	Describe("FindUsersByIDs", func() {
		It("returns the users with the given ids", func() {
			users, err := storage.FindUsersByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(BeEmpty())

			err = db.Create(&User{}).Error
			Expect(err).ToNot(HaveOccurred())

			users, err = storage.FindUsersByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(HaveLen(1))
		})
	})

	Describe("FindUserByEmail", func() {
		It("returns the user with the given email", func() {
			user, err := storage.FindUserByEmail("user@systemli.org")
			Expect(err).To(HaveOccurred())
			Expect(user).To(BeZero())

			err = db.Create(&User{Email: "user@systemli.org"}).Error
			Expect(err).ToNot(HaveOccurred())

			user, err = storage.FindUserByEmail("user@systemli.org")
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeZero())
		})
	})

	Describe("FindUsersByTicker", func() {
		It("returns the users with the given ticker", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			users, err := storage.FindUsersByTicker(ticker)
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(BeEmpty())

			user, err := NewUser("user@systemli.org", "password")
			Expect(err).ToNot(HaveOccurred())
			err = storage.SaveUser(&user)
			Expect(err).ToNot(HaveOccurred())

			ticker.Users = append(ticker.Users, user)
			err = storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			users, err = storage.FindUsersByTicker(ticker)
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(HaveLen(1))
		})
	})

	Describe("SaveUser", func() {
		It("persists the user", func() {
			user, err := NewUser("user@systemli.org", "password")
			Expect(err).ToNot(HaveOccurred())

			err = storage.SaveUser(&user)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&User{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})

	Describe("DeleteUser", func() {
		It("deletes the user", func() {
			user, err := NewUser("user@systemli.org", "password")
			Expect(err).ToNot(HaveOccurred())

			err = storage.SaveUser(&user)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&User{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteUser(user)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&User{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("DeleteTickerUser", func() {
		It("deletes the user from the ticker", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			user, err := NewUser("user@systemli.org", "password")
			Expect(err).ToNot(HaveOccurred())
			err = storage.SaveUser(&user)
			Expect(err).ToNot(HaveOccurred())

			ticker.Users = append(ticker.Users, user)
			err = storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&User{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteTickerUser(&ticker, &user)
			Expect(err).ToNot(HaveOccurred())
			Expect(ticker.Users).To(BeEmpty())
		})
	})

	Describe("AddTickerUser", func() {
		It("adds the user to the ticker", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			user, err := NewUser("user@systemli.org", "password")
			Expect(err).ToNot(HaveOccurred())
			err = storage.SaveUser(&user)
			Expect(err).ToNot(HaveOccurred())

			err = storage.AddTickerUser(&ticker, &user)
			Expect(err).ToNot(HaveOccurred())
			Expect(ticker.Users).To(HaveLen(1))
		})
	})

	Describe("FindTickers", func() {
		It("returns all tickers", func() {
			tickers, err := storage.FindTickers()
			Expect(err).ToNot(HaveOccurred())
			Expect(tickers).To(BeEmpty())

			err = db.Create(&Ticker{}).Error
			Expect(err).ToNot(HaveOccurred())

			tickers, err = storage.FindTickers()
			Expect(err).ToNot(HaveOccurred())
			Expect(tickers).To(HaveLen(1))
		})
	})

	Describe("FindTickersByIDs", func() {
		It("returns the tickers with the given ids", func() {
			tickers, err := storage.FindTickersByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(tickers).To(BeEmpty())

			err = db.Create(&Ticker{}).Error
			Expect(err).ToNot(HaveOccurred())

			tickers, err = storage.FindTickersByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(tickers).To(HaveLen(1))
		})
	})

	Describe("FindTickerByID", func() {
		It("returns the ticker with the given id", func() {
			ticker, err := storage.FindTickerByID(1)
			Expect(err).To(HaveOccurred())
			Expect(ticker).To(BeZero())

			err = db.Create(&Ticker{}).Error
			Expect(err).ToNot(HaveOccurred())

			ticker, err = storage.FindTickerByID(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(ticker).ToNot(BeZero())
		})
	})

	Describe("FindTickerByDomain", func() {
		It("returns the ticker with the given domain", func() {
			ticker, err := storage.FindTickerByDomain("systemli.org")
			Expect(err).To(HaveOccurred())
			Expect(ticker).To(BeZero())

			err = db.Create(&Ticker{Domain: "systemli.org"}).Error
			Expect(err).ToNot(HaveOccurred())

			ticker, err = storage.FindTickerByDomain("systemli.org")
			Expect(err).ToNot(HaveOccurred())
			Expect(ticker).ToNot(BeZero())
		})
	})

	Describe("SaveTicker", func() {
		It("persists the ticker", func() {
			ticker := NewTicker()

			err = storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Ticker{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})

	Describe("DeleteTicker", func() {
		It("deletes the ticker", func() {
			ticker := NewTicker()

			err = storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Ticker{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteTicker(ticker)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&Ticker{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("FindUploadByUUID", func() {
		It("returns the upload with the given uuid", func() {
			upload, err := storage.FindUploadByUUID("uuid")
			Expect(err).To(HaveOccurred())
			Expect(upload).To(BeZero())

			err = db.Create(&Upload{UUID: "uuid"}).Error
			Expect(err).ToNot(HaveOccurred())

			upload, err = storage.FindUploadByUUID("uuid")
			Expect(err).ToNot(HaveOccurred())
			Expect(upload).ToNot(BeZero())
		})
	})

	Describe("FindUploadsByIDs", func() {
		It("returns the uploads with the given ids", func() {
			uploads, err := storage.FindUploadsByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(uploads).To(BeEmpty())

			err = db.Create(&Upload{}).Error
			Expect(err).ToNot(HaveOccurred())

			uploads, err = storage.FindUploadsByIDs([]int{1, 2})
			Expect(err).ToNot(HaveOccurred())
			Expect(uploads).To(HaveLen(1))
		})
	})

	Describe("SaveUpload", func() {
		It("persists the upload", func() {
			upload := NewUpload("image.jpg", "content-type", 1)

			err = storage.SaveUpload(&upload)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})

	Describe("DeleteUpload", func() {
		It("deletes the upload", func() {
			upload := NewUpload("image.jpg", "content-type", 1)

			err = storage.SaveUpload(&upload)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteUpload(upload)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("DeleteUploads", func() {
		It("deletes the uploads", func() {
			upload := NewUpload("image.jpg", "content-type", 1)

			err = storage.SaveUpload(&upload)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			uploads := []Upload{upload}
			storage.DeleteUploads(uploads)

			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("DeleteUploadsByTicker", func() {
		It("deletes the uploads", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			upload := NewUpload("image.jpg", "content-type", ticker.ID)
			err = storage.SaveUpload(&upload)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteUploadsByTicker(ticker)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&Upload{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("FindMessage", func() {
		It("returns the message with the given id", func() {
			message, err := storage.FindMessage(1, 1)
			Expect(err).To(HaveOccurred())
			Expect(message).To(BeZero())

			err = db.Create(&Message{ID: 1, TickerID: 1}).Error
			Expect(err).ToNot(HaveOccurred())

			message, err = storage.FindMessage(1, 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(message).ToNot(BeZero())
		})
	})

	Describe("FindMessagesByTicker", func() {
		It("returns the messages with the given ticker", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			messages, err := storage.FindMessagesByTicker(ticker)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(BeEmpty())

			err = db.Create(&Message{TickerID: ticker.ID}).Error
			Expect(err).ToNot(HaveOccurred())

			messages, err = storage.FindMessagesByTicker(ticker)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(1))
		})
	})

	Describe("FindMessagesByTickerAndPagination", func() {
		It("returns the messages with the given ticker and pagination", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			c := &gin.Context{}
			p := pagination.NewPagination(c)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *p)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(BeEmpty())

			err = db.Create(&Message{TickerID: ticker.ID}).Error
			Expect(err).ToNot(HaveOccurred())

			messages, err = storage.FindMessagesByTickerAndPagination(ticker, *p)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(1))

			err = db.Create([]Message{
				{TickerID: ticker.ID, ID: 2},
				{TickerID: ticker.ID, ID: 3},
				{TickerID: ticker.ID, ID: 4},
			}).Error
			Expect(err).ToNot(HaveOccurred())

			c = &gin.Context{}
			c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2"}}
			p = pagination.NewPagination(c)
			messages, err = storage.FindMessagesByTickerAndPagination(ticker, *p)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(2))

			c = &gin.Context{}
			c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2&after=2"}}
			p = pagination.NewPagination(c)
			messages, err = storage.FindMessagesByTickerAndPagination(ticker, *p)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(2))

			c = &gin.Context{}
			c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2&before=4"}}
			p = pagination.NewPagination(c)
			messages, err = storage.FindMessagesByTickerAndPagination(ticker, *p)
			Expect(err).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(2))
		})
	})

	Describe("SaveMessage", func() {
		It("persists the message", func() {
			message := NewMessage()

			err = storage.SaveMessage(&message)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Message{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
		})
	})

	Describe("DeleteMessage", func() {
		It("deletes the message", func() {
			message := NewMessage()

			err = storage.SaveMessage(&message)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Message{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteMessage(message)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&Message{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("DeleteMessages", func() {
		It("deletes the messages", func() {
			ticker := NewTicker()
			err := storage.SaveTicker(&ticker)
			Expect(err).ToNot(HaveOccurred())

			message := NewMessage()
			message.TickerID = ticker.ID
			err = storage.SaveMessage(&message)
			Expect(err).ToNot(HaveOccurred())

			var count int64
			err = db.Model(&Message{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(1)))

			err = storage.DeleteMessages(ticker)
			Expect(err).ToNot(HaveOccurred())

			err = db.Model(&Message{}).Count(&count).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("GetInactiveSettings", func() {
		It("returns the default inactive setting", func() {
			setting := storage.GetInactiveSettings()
			Expect(setting.Author).To(Equal(DefaultInactiveSettings().Author))
		})

		It("returns the inactive setting", func() {
			settings := InactiveSettings{
				Author: "author",
			}

			err = storage.SaveInactiveSettings(settings)
			Expect(err).ToNot(HaveOccurred())

			setting := storage.GetInactiveSettings()
			Expect(setting.Author).To(Equal(settings.Author))
		})
	})

	Describe("GetRefreshIntervalSetting", func() {
		It("returns the default refresh interval setting", func() {
			setting := storage.GetRefreshIntervalSettings()
			Expect(setting.RefreshInterval).To(Equal(DefaultRefreshIntervalSettings().RefreshInterval))
		})

		It("returns the refresh interval setting", func() {
			settings := RefreshIntervalSettings{
				RefreshInterval: 1000,
			}

			err = storage.SaveRefreshIntervalSettings(settings)
			Expect(err).ToNot(HaveOccurred())

			setting := storage.GetRefreshIntervalSettings()
			Expect(setting.RefreshInterval).To(Equal(settings.RefreshInterval))
		})
	})
})
