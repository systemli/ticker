package legacy

import (
	"os"
	"testing"

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

		ticker     Ticker
		message    Message
		attachment Attachment
		adminUser  User
		user       User
		upload     Upload
	)

	ticker = Ticker{
		ID:          161,
		Active:      true,
		Title:       "Ticker",
		Domain:      "ticker.org",
		Description: "Ticker description",
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
		ID:          1,
		Ticker:      161,
		Text:        "Message",
		Attachments: []Attachment{attachment},
	}

	adminUser = User{
		Email:        "admin@systemli.org",
		IsSuperAdmin: true,
	}

	user = User{
		Email:        "user@systemli.org",
		IsSuperAdmin: false,
		Tickers:      []int{161},
	}

	upload = Upload{
		UUID:        "uuid",
		TickerID:    161,
		Extension:   "png",
		ContentType: "image/png",
		Path:        "testdata/uploads",
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

			var telegram storage.TickerTelegram
			Expect(newDb.First(&telegram).Error).ToNot(HaveOccurred())
			Expect(telegram.Active).To(BeTrue())

			var mastodon storage.TickerMastodon
			Expect(newDb.First(&mastodon).Error).ToNot(HaveOccurred())
			Expect(mastodon.Active).To(BeTrue())

			var users []storage.User
			Expect(newDb.Find(&users).Error).ToNot(HaveOccurred())
			Expect(users).To(HaveLen(2))
			Expect(users[0].Email).To(Equal(adminUser.Email))
			Expect(users[0].IsSuperAdmin).To(BeTrue())
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

			var attachments []storage.Attachment
			Expect(newDb.Find(&attachments).Error).ToNot(HaveOccurred())
			Expect(attachments).To(HaveLen(1))
			Expect(attachments[0].UUID).To(Equal(attachment.UUID))

			var uploads []storage.Upload
			Expect(newDb.Find(&uploads).Error).ToNot(HaveOccurred())
			Expect(uploads).To(HaveLen(1))
			Expect(uploads[0].UUID).To(Equal(upload.UUID))
		})
	})
})
