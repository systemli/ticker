package storage

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type MigrationTestSuite struct {
	db *gorm.DB
	suite.Suite
}

func (s *MigrationTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open("file:testdatabase?mode=memory&cache=shared"), &gorm.Config{})
	s.NoError(err)

	err = db.AutoMigrate(
		&Ticker{},
		&TickerTelegram{},
		&TickerMastodon{},
		&TickerBluesky{},
		&TickerSignalGroup{},
		&TickerWebsite{},
		&User{},
		&Message{},
		&Upload{},
		&Attachment{},
		&Setting{},
	)
	s.NoError(err)

	s.db = db
}

func (s *MigrationTestSuite) TestMigrateDB() {
	s.Run("without existing data", func() {
		err := MigrateDB(s.db)
		s.NoError(err)
	})

	s.Run("with existing data", func() {
		ticker := Ticker{Domain: "example.org"}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		err = MigrateDB(s.db)
		s.NoError(err)

		var tickerWebsite TickerWebsite
		err = s.db.First(&tickerWebsite).Error
		s.NoError(err)
		s.Equal(tickerWebsite.TickerID, ticker.ID)
		s.Equal(tickerWebsite.Origin, "https://example.org")
	})
}

func (s *MigrationTestSuite) TestMigrateDomain() {
	s.Run("when domain is localhost", func() {
		s.Equal("http://localhost", migrateDomain("localhost"))
	})

	s.Run("when domain is not localhost", func() {
		s.Equal("https://example.org", migrateDomain("example.org"))
	})
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
