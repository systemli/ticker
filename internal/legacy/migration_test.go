package legacy

import (
	"os"
	"testing"
	"time"

	"github.com/asdine/storm"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MigrationTestSuite struct {
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

	suite.Suite
}

func (s *MigrationTestSuite) SetupTest() {
	s.oldDb, s.err = storm.Open("storm.db", storm.BoltOptions(0600, nil))
	s.NoError(s.err)

	s.newDb, s.err = gorm.Open(sqlite.Open("file:testdatabase?mode=memory&cache=shared"), &gorm.Config{})
	s.NoError(s.err)
	s.NoError(storage.MigrateDB(s.newDb))

	oldStorage := NewLegacyStorage(s.oldDb)
	newStorage := storage.NewSqlStorage(s.newDb, "testdata/uploads")
	s.migration = NewMigration(oldStorage, newStorage)

	s.ticker = Ticker{
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

	s.attachment = Attachment{
		UUID:        "uuid",
		Extension:   "png",
		ContentType: "image/png",
	}

	s.message = Message{
		ID:           1,
		CreationDate: time.Now(),
		Ticker:       161,
		Text:         "Message",
		Attachments:  []Attachment{s.attachment},
	}

	s.adminUser = User{
		ID:                1,
		Email:             "admin@systemli.org",
		EncryptedPassword: "notempty",
		CreationDate:      time.Now(),
		IsSuperAdmin:      true,
	}

	s.user = User{
		ID:                2,
		Email:             "user@systemli.org",
		EncryptedPassword: "notempty",
		CreationDate:      time.Now(),
		IsSuperAdmin:      false,
		Tickers:           []int{161},
	}

	s.upload = Upload{
		ID:           1,
		CreationDate: time.Now(),
		TickerID:     161,
		UUID:         "uuid",
		Extension:    "png",
		ContentType:  "image/png",
		Path:         "testdata/uploads",
	}

	s.refreshIntervalSetting = Setting{
		ID:    1,
		Name:  "refresh_interval",
		Value: 10000,
	}

	s.inactiveSettingsSetting = Setting{
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

	s.NoError(s.oldDb.Save(&s.ticker))
	s.NoError(s.oldDb.Save(&s.message))
	s.NoError(s.oldDb.Save(&s.adminUser))
	s.NoError(s.oldDb.Save(&s.user))
	s.NoError(s.oldDb.Save(&s.upload))
	s.NoError(s.oldDb.Save(&s.refreshIntervalSetting))
	s.NoError(s.oldDb.Save(&s.inactiveSettingsSetting))
}

func (s *MigrationTestSuite) TearDownTest() {
	s.NoError(s.oldDb.Close())
	s.NoError(os.Remove("storm.db"))
}

func (s *MigrationTestSuite) TestDo() {
	s.err = s.migration.Do()
	s.NoError(s.err)

	var tickers []storage.Ticker
	s.NoError(s.newDb.Find(&tickers).Error)
	s.Equal(1, len(tickers))
	s.Equal(s.ticker.ID, tickers[0].ID)
	s.Equal(s.ticker.Active, tickers[0].Active)

	var telegram storage.TickerTelegram
	s.NoError(s.newDb.First(&telegram).Error)
	s.Equal(s.ticker.Telegram.Active, telegram.Active)

	var mastodon storage.TickerMastodon
	s.NoError(s.newDb.First(&mastodon).Error)
	s.Equal(s.ticker.Mastodon.Active, mastodon.Active)

	var users []storage.User
	s.NoError(s.newDb.Find(&users).Error)
	s.Equal(2, len(users))
	s.Equal(s.adminUser.ID, users[0].ID)
	s.Equal(s.adminUser.Email, users[0].Email)
	s.Equal(s.adminUser.IsSuperAdmin, users[0].IsSuperAdmin)
	s.Equal(s.user.ID, users[1].ID)
	s.Equal(s.user.Email, users[1].Email)
	s.Equal(s.user.IsSuperAdmin, users[1].IsSuperAdmin)

	var tickersUsers []storage.User
	s.NoError(s.newDb.Model(&tickers[0]).Association("Users").Find(&tickersUsers))
	s.Equal(1, len(tickersUsers))
	s.Equal(s.user.ID, tickersUsers[0].ID)
	s.Equal(s.user.Email, tickersUsers[0].Email)

	var messages []storage.Message
	s.NoError(s.newDb.Find(&messages).Error)
	s.Equal(1, len(messages))
	s.Equal(s.message.ID, messages[0].ID)
	s.Equal(s.message.Ticker, messages[0].TickerID)
	s.Equal(s.message.Text, messages[0].Text)

	var attachments []storage.Attachment
	s.NoError(s.newDb.Find(&attachments).Error)
	s.Equal(1, len(attachments))
	s.Equal(s.attachment.UUID, attachments[0].UUID)
	s.Equal(s.attachment.Extension, attachments[0].Extension)
	s.Equal(s.attachment.ContentType, attachments[0].ContentType)

	var uploads []storage.Upload
	s.NoError(s.newDb.Find(&uploads).Error)
	s.Equal(1, len(uploads))
	s.Equal(s.upload.ID, uploads[0].ID)
	s.Equal(s.upload.TickerID, uploads[0].TickerID)
	s.Equal(s.upload.UUID, uploads[0].UUID)
	s.Equal(s.upload.Extension, uploads[0].Extension)
	s.Equal(s.upload.ContentType, uploads[0].ContentType)
	s.Equal(s.upload.Path, uploads[0].Path)

	var settings []storage.Setting
	s.NoError(s.newDb.Find(&settings).Error)
	s.Equal(2, len(settings))
	s.Equal(s.refreshIntervalSetting.Name, settings[0].Name)
	s.Equal(`{"refreshInterval":10000}`, settings[0].Value)
	s.Equal(s.inactiveSettingsSetting.Name, settings[1].Name)
	s.Equal(`{"headline":"Headline","subHeadline":"Subheadline","description":"Description","author":"Author","email":"Email","homepage":"","twitter":"Twitter"}`, settings[1].Value)
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
