package storage

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/systemli/ticker/internal/api/pagination"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = Describe("Storage", func() {
	var storagePath = fmt.Sprintf("%s/storage_%d.db", strings.TrimSuffix(os.TempDir(), "/"), time.Now().Unix())
	var storage = NewStorage(storagePath, "/uploads")

	BeforeEach(func() {
		storage.DropAll()
	})

	When("no ticker is present", func() {
		It("shouldn't find any ticker", func() {
			tickers, err := storage.FindTickers()
			Expect(err).To(BeNil())
			Expect(tickers).To(HaveLen(0))
		})
	})

	When("one ticker is present", func() {
		domain := "ticker.systemli.org"
		ticker := Ticker{CreationDate: time.Now(), Domain: domain}

		BeforeEach(func() {
			err := storage.SaveTicker(&ticker)
			Expect(err).To(BeNil())
		})

		It("should find one ticker", func() {
			tickers, err := storage.FindTickers()
			Expect(err).To(BeNil())
			Expect(tickers).To(HaveLen(1))
		})

		It("should find one ticker by domain", func() {
			ticker, err := storage.FindTickerByDomain(domain)
			Expect(err).To(BeNil())
			Expect(ticker).NotTo(BeNil())
		})

		It("should find one ticker by id", func() {
			ticker, err := storage.FindTickerByID(1)
			Expect(err).To(BeNil())
			Expect(ticker).NotTo(BeNil())
		})

		It("shouldn't find one ticker by misspelled domain", func() {
			_, err := storage.FindTickerByDomain("missspelled")
			Expect(err).NotTo(BeNil())
		})

		It("shouldn't find one ticker by nonexisting id", func() {
			_, err := storage.FindTickerByID(100)
			Expect(err).NotTo(BeNil())
		})

		It("should be possible to delete them", func() {
			ticker, _ := storage.FindTickerByID(1)
			err := storage.DeleteTicker(ticker)
			Expect(err).To(BeNil())
		})

		It("should be possible to add users", func() {
			user, _ := NewUser("louis@systemli.org", "password")
			ticker, _ := storage.FindTickerByID(1)
			err := storage.SaveUser(&user)
			Expect(err).To(BeNil())

			err = storage.AddUsersToTicker(ticker, []int{user.ID})
			Expect(err).To(BeNil())
			user, _ = storage.FindUserByEmail("louis@systemli.org")
			Expect(user.Tickers).To(HaveLen(1))
		})

		It("should not add ticker to admin users", func() {
			user, _ := NewAdminUser("louis@systemli.org", "password")
			ticker, _ := storage.FindTickerByID(1)
			err := storage.SaveUser(&user)
			Expect(err).To(BeNil())

			err = storage.AddUsersToTicker(ticker, []int{user.ID})
			Expect(err).To(BeNil())
			user, _ = storage.FindUserByEmail("louis@systemli.org")
			Expect(user.Tickers).To(HaveLen(0))
		})

		It("should be possible to remove users", func() {
			ticker, _ := storage.FindTickerByID(1)
			user, _ := NewUser("louis@systemli.org", "password")
			user.AddTicker(ticker)
			err := storage.SaveUser(&user)
			Expect(err).To(BeNil())

			err = storage.RemoveTickerFromUser(ticker, user)
			Expect(err).To(BeNil())
			user, _ = storage.FindUserByEmail("louis@systemli.org")
			Expect(user.Tickers).To(HaveLen(0))

		})
	})

	When("two tickers are present", func() {
		domain1 := "ticker.systemli.org"
		domain2 := "ticker.tem.li"
		ticker1 := Ticker{CreationDate: time.Now(), Domain: domain1}
		ticker2 := Ticker{CreationDate: time.Now(), Domain: domain2}

		BeforeEach(func() {
			err := storage.SaveTicker(&ticker1)
			Expect(err).To(BeNil())
			err = storage.SaveTicker(&ticker2)
			Expect(err).To(BeNil())
		})

		It("should find two ticker", func() {
			tickers, err := storage.FindTickers()
			Expect(err).To(BeNil())
			Expect(tickers).To(HaveLen(2))
		})

		It("should find two tickers by their ids", func() {
			tickers, err := storage.FindTickersByIDs([]int{ticker1.ID, ticker2.ID})
			Expect(err).To(BeNil())
			Expect(tickers).To(HaveLen(2))
		})
	})

	When("one user is present", func() {
		email := "user@systemli.org"
		user := User{Email: email, Tickers: []int{2}}

		BeforeEach(func() {
			err := storage.SaveUser(&user)
			Expect(err).To(BeNil())
		})

		It("should find one user", func() {
			users, err := storage.FindUsers()
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(1))
		})

		It("should find one user by email", func() {
			user, err := storage.FindUserByEmail(email)
			Expect(err).To(BeNil())
			Expect(user.Email).To(Equal(email))
		})

		It("should find one user by ID", func() {
			user, err := storage.FindUserByID(1)
			Expect(err).To(BeNil())
			Expect(user.ID).To(Equal(1))
		})

		It("should be possible to delete them", func() {
			user, _ := storage.FindUserByID(1)
			err := storage.DeleteUser(user)
			Expect(err).To(BeNil())
		})

		It("should find users for a ticker with id 2", func() {
			ticker := Ticker{ID: 2}
			users, err := storage.FindUsersByTicker(ticker)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(1))
		})

		It("shouldn't find users for a ticker with id 1", func() {
			ticker := Ticker{ID: 1}
			users, err := storage.FindUsersByTicker(ticker)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(0))
		})

		It("should be return the correct count", func() {
			c, err := storage.CountUser()
			Expect(err).To(BeNil())
			Expect(c).To(Equal(1))
		})
	})

	When("no messages are present", func() {
		It("shouldn't find any messages", func() {
			ticker := Ticker{ID: 1, Active: true}
			ctx := gin.Context{}
			pagination := pagination.NewPagination(&ctx)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *pagination)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(0))
		})
	})

	When("messages are present", func() {
		var ticker Ticker
		var message Message

		BeforeEach(func() {
			ticker = NewTicker()
			ticker.Active = true
			err := storage.SaveTicker(&ticker)
			Expect(err).To(BeNil())
			message = NewMessage()
			message.Ticker = ticker.ID
			err = storage.SaveMessage(&message)
			Expect(err).To(BeNil())
		})

		It("should return all messages", func() {
			messages, err := storage.FindMessagesByTicker(ticker)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(1))
		})

		It("should return no messages", func() {
			ticker1 := NewTicker()
			messages, err := storage.FindMessagesByTicker(ticker1)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(0))
		})

		It("should return messages on active ticker", func() {
			ctx := gin.Context{}
			pagination := pagination.NewPagination(&ctx)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *pagination)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(1))
		})

		It("should not return newer messages on active ticker", func() {
			ctx := gin.Context{
				Request: &http.Request{
					URL: &url.URL{RawQuery: fmt.Sprintf("after=%d", message.ID)},
				},
			}
			pagination := pagination.NewPagination(&ctx)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *pagination)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(0))
		})

		It("should not return older messages on active ticker", func() {
			ctx := gin.Context{
				Request: &http.Request{
					URL: &url.URL{RawQuery: fmt.Sprintf("before=%d", message.ID)},
				},
			}
			pagination := pagination.NewPagination(&ctx)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *pagination)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(0))
		})

		It("should not return messages for inactive ticker", func() {
			ticker.Active = false
			err := storage.SaveTicker(&ticker)
			Expect(err).To(BeNil())

			ctx := gin.Context{}
			pagination := pagination.NewPagination(&ctx)
			messages, err := storage.FindMessagesByTickerAndPagination(ticker, *pagination)
			Expect(err).To(BeNil())
			Expect(messages).To(HaveLen(0))
		})

		It("should return the message when queried", func() {
			found, err := storage.FindMessage(message.Ticker, message.ID)
			Expect(err).To(BeNil())
			Expect(found.ID).To(Equal(message.ID))
		})

		It("should return no uploads", func() {
			uploads := storage.FindUploadsByMessage(message)
			Expect(uploads).To(HaveLen(0))
		})

		It("should return no uploads when reference is invalid", func() {
			message.Attachments = []Attachment{{UUID: "invalid", ContentType: "text/plain", Extension: "txt"}}
			_ = storage.SaveMessage(&message)

			uploads := storage.FindUploadsByMessage(message)
			Expect(uploads).To(HaveLen(0))
		})

		It("should return uploads", func() {
			upload := NewUpload("image.jpg", "image/jpeg", 1)
			err := storage.SaveUpload(&upload)
			Expect(err).To(BeNil())
			message.AddAttachment(upload)
			err = storage.SaveMessage(&message)
			Expect(err).To(BeNil())

			uploads := storage.FindUploadsByMessage(message)
			Expect(uploads).To(HaveLen(1))
		})

		It("should be possible to delete the message", func() {
			err := storage.DeleteMessage(message)
			Expect(err).To(BeNil())
		})

		It("should be possible to delete messages by ticker", func() {
			err := storage.DeleteMessages(ticker)
			Expect(err).To(BeNil())
		})
	})

	When("uploads are present", func() {
		ticker := NewTicker()
		_ = storage.SaveTicker(&ticker)
		upload := NewUpload("filename.txt", "text/plain", ticker.ID)
		err := storage.SaveUpload(&upload)
		Expect(err).To(BeNil())

		It("should return upload by uuid", func() {
			u, err := storage.FindUploadByUUID(upload.UUID)
			Expect(err).To(BeNil())
			Expect(u.UUID).To(Equal(upload.UUID))
		})

		It("should return upload by ids", func() {
			u, err := storage.FindUploadsByIDs([]int{upload.ID})
			Expect(err).To(BeNil())
			Expect(u).To(HaveLen(1))
		})

		It("can deleted a specific one", func() {
			err := storage.DeleteUpload(upload)
			Expect(err).To(BeNil())
		})

		It("can delete multiple ones", func() {
			storage.DeleteUploads([]Upload{upload})
			_, err := storage.FindUploadByUUID(upload.UUID)
			Expect(err).NotTo(BeNil())
		})

		It("can be deleted by ticker", func() {
			err := storage.DeleteUploadsByTicker(ticker)
			Expect(err).To(BeNil())
			_, err = storage.FindUploadByUUID(upload.UUID)
			Expect(err).NotTo(BeNil())
		})
	})

	When("settings are fetched", func() {
		It("should return default refresh interval", func() {
			refreshInterval := storage.GetRefreshIntervalSetting()
			Expect(refreshInterval.Value.(float64)).To(Equal(SettingDefaultRefreshInterval))
		})

		It("should return default inactive settings", func() {
			InactiveSetting := storage.GetInactiveSetting()
			Expect(InactiveSetting.Value).To(Equal(DefaultInactiveSettings()))
		})

		refreshInterval := float64(20000)
		err := storage.SaveRefreshInterval(refreshInterval)
		fmt.Println(err)
		Expect(err).To(BeNil())

		It("should return the user defined refresh interval", func() {
			err := storage.SaveRefreshInterval(2000)
			Expect(err).To(BeNil())
			refreshIntervalSetting := storage.GetRefreshIntervalSetting()
			Expect(refreshIntervalSetting.Value).To(Equal(float64(2000)))
			refreshInterval := storage.GetRefreshIntervalSettingValue()
			Expect(refreshInterval).To(Equal(2000))
		})

		It("should update the existing refresh interval", func() {
			err := storage.SaveRefreshInterval(2000)
			Expect(err).To(BeNil())
			err = storage.SaveRefreshInterval(1000)
			Expect(err).To(BeNil())
		})

		It("should return the user defined inactive settings", func() {
			inactiveSettings := InactiveSettings{
				Headline:    "Headline",
				SubHeadline: "SubHeadline",
				Description: "New Description",
				Author:      "Author",
				Email:       "Email",
				Homepage:    "Homepage",
				Twitter:     "Twitter",
			}
			err := storage.SaveInactiveSetting(inactiveSettings)
			Expect(err).To(BeNil())
			setting := storage.GetInactiveSetting()
			Expect(setting.Name).To(Equal(SettingInactiveName))
		})
	})
})
