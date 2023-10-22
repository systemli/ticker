package legacy

import (
	"os"
	"testing"
	"time"

	"github.com/asdine/storm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/systemli/ticker/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Migration Suite")
}

var _ = Describe("Migration", func() {
	var (
		oldDb     *storm.DB
		newDb     *gorm.DB
		migration *Migration
		err       error

		ticker                  Ticker
		message                 Message
		attachment              Attachment
		adminUser               User
		user                    User
		upload                  Upload
		refreshIntervalSetting  Setting
		inactiveSettingsSetting Setting
	)

	ticker = Ticker{
		ID:           161,
		CreationDate: time.Now(),
		Active:       true,
		Title:        "Ticker",
		Domain:       "ticker.org",
		Description:  "Ticker description",
		Information: Information{
			Author: "Author",
		},
		Telegram: Telegram{
			Active: true,
		},
		Mastodon: Mastodon{
			Active: true,
		},
	}

	attachment = Attachment{
		UUID:        "uuid",
		Extension:   "png",
		ContentType: "image/png",
	}

	message = Message{
		ID:           1,
		CreationDate: time.Now(),
		Ticker:       161,
		Text:         "Message",
		Attachments:  []Attachment{attachment},
	}

	adminUser = User{
		ID:           1,
		Email:        "admin@systemli.org",
		CreationDate: time.Now(),
		IsSuperAdmin: true,
	}

	user = User{
		ID:           2,
		Email:        "user@systemli.org",
		CreationDate: time.Now(),
		IsSuperAdmin: false,
		Tickers:      []int{161},
	}

	upload = Upload{
		ID:           1,
		CreationDate: time.Now(),
		TickerID:     161,
		UUID:         "uuid",
		Extension:    "png",
		ContentType:  "image/png",
		Path:         "testdata/uploads",
	}

	refreshIntervalSetting = Setting{
		ID:    1,
		Name:  "refresh_interval",
		Value: 10000,
	}

	inactiveSettingsSetting = Setting{
		ID:   2,
		Name: "inactive_settings",
		Value: map[string]string{
			"headline":     "Headline",
			"sub_headline": "Subheadline",
			"description":  "Description",
			"author":       "Author",
			"email":        "Email",
			"twitter":      "Twitter",
		},
	}

	BeforeEach(func() {
		oldDb, err = storm.Open("storm.db", storm.BoltOptions(0600, nil))
		Expect(err).ToNot(HaveOccurred())

		newDb, err = gorm.Open(sqlite.Open("file:testdatabase?mode=memory&cache=shared"), &gorm.Config{})
		Expect(err).ToNot(HaveOccurred())
		Expect(storage.MigrateDB(newDb)).To(Succeed())

		oldStorage := NewLegacyStorage(oldDb)
		newStorage := storage.NewSqlStorage(newDb, "testdata/uploads")
		migration = NewMigration(oldStorage, newStorage)

		Expect(oldDb.Save(&ticker)).To(Succeed())
		Expect(oldDb.Save(&message)).To(Succeed())
		Expect(oldDb.Save(&adminUser)).To(Succeed())
		Expect(oldDb.Save(&user)).To(Succeed())
		Expect(oldDb.Save(&upload)).To(Succeed())
		Expect(oldDb.Save(&refreshIntervalSetting)).To(Succeed())
		Expect(oldDb.Save(&inactiveSettingsSetting)).To(Succeed())
	})

	AfterEach(func() {
		Expect(oldDb.Close()).To(Succeed())
		Expect(os.Remove("storm.db")).To(Succeed())
	})

	Describe("Do", func() {
		It("migrates all the data successfully", func() {
			err = migration.Do()
			Expect(err).ToNot(HaveOccurred())

			var tickers []storage.Ticker
			Expect(newDb.Find(&tickers).Error).ToNot(HaveOccurred())
			Expect(tickers).To(HaveLen(1))
			Expect(tickers[0].ID).To(Equal(ticker.ID))
			Expect(tickers[0].CreatedAt).Should(BeTemporally("~", ticker.CreationDate, time.Second))
			Expect(tickers[0].UpdatedAt).Should(BeTemporally("~", ticker.CreationDate, time.Second))
			Expect(tickers[0].Active).To(BeTrue())

			var telegram storage.TickerTelegram
			Expect(newDb.First(&telegram).Error).ToNot(HaveOccurred())
			Expect(telegram.Active).To(BeTrue())

			var mastodon storage.TickerMastodon
			Expect(newDb.First(&mastodon).Error).ToNot(HaveOccurred())
			Expect(mastodon.Active).To(BeTrue())

			var users []storage.User
			Expect(newDb.Find(&users).Error).ToNot(HaveOccurred())
			Expect(users).To(HaveLen(2))
			Expect(users[0].ID).To(Equal(adminUser.ID))
			Expect(users[0].CreatedAt).Should(BeTemporally("~", adminUser.CreationDate, time.Second))
			Expect(users[0].Email).To(Equal(adminUser.Email))
			Expect(users[0].IsSuperAdmin).To(BeTrue())
			Expect(users[1].ID).To(Equal(user.ID))
			Expect(users[1].CreatedAt).Should(BeTemporally("~", user.CreationDate, time.Second))
			Expect(users[1].Email).To(Equal(user.Email))
			Expect(users[1].IsSuperAdmin).To(BeFalse())

			var tickersUsers []storage.User
			Expect(newDb.Model(&tickers[0]).Association("Users").Find(&tickersUsers)).To(Succeed())
			Expect(tickersUsers).To(HaveLen(1))
			Expect(tickersUsers[0].Email).To(Equal(user.Email))

			var messages []storage.Message
			Expect(newDb.Find(&messages).Error).ToNot(HaveOccurred())
			Expect(messages).To(HaveLen(1))
			Expect(messages[0].ID).To(Equal(message.ID))
			Expect(messages[0].TickerID).To(Equal(message.Ticker))
			Expect(messages[0].Text).To(Equal(message.Text))
			Expect(messages[0].CreatedAt).Should(BeTemporally("~", message.CreationDate, time.Second))

			var attachments []storage.Attachment
			Expect(newDb.Find(&attachments).Error).ToNot(HaveOccurred())
			Expect(attachments).To(HaveLen(1))
			Expect(attachments[0].MessageID).To(Equal(message.ID))
			Expect(attachments[0].UUID).To(Equal(attachment.UUID))
			Expect(attachments[0].CreatedAt).Should(BeTemporally("~", message.CreationDate, time.Second))

			var uploads []storage.Upload
			Expect(newDb.Find(&uploads).Error).ToNot(HaveOccurred())
			Expect(uploads).To(HaveLen(1))
			Expect(uploads[0].ID).To(Equal(upload.ID))
			Expect(uploads[0].CreatedAt).Should(BeTemporally("~", upload.CreationDate, time.Second))
			Expect(uploads[0].TickerID).To(Equal(upload.TickerID))
			Expect(uploads[0].UUID).To(Equal(upload.UUID))

			var settings []storage.Setting
			Expect(newDb.Find(&settings).Error).ToNot(HaveOccurred())
			Expect(settings).To(HaveLen(2))
			Expect(settings[0].Name).To(Equal(refreshIntervalSetting.Name))
			Expect(settings[0].Value).To(Equal(`{"refreshInterval":10000}`))
			Expect(settings[1].Name).To(Equal(inactiveSettingsSetting.Name))
			Expect(settings[1].Value).To(Equal(`{"headline":"Headline","subHeadline":"Subheadline","description":"Description","author":"Author","email":"Email","homepage":"","twitter":"Twitter"}`))
		})
	})
})
