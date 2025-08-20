package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
}

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
