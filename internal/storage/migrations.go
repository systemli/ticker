package storage

import (
	"gorm.io/gorm"
)

// MigrateDB migrates the database
func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&Ticker{},
		&TickerMastodon{},
		&TickerTelegram{},
		&TickerBluesky{},
		&TickerSignalGroup{},
		&TickerWebsite{},
		&User{},
		&Setting{},
		&Upload{},
		&Message{},
		&Attachment{},
	); err != nil {
		return err
	}

	// Drop the column geo_information from Message if it exists
	if db.Migrator().HasColumn(&Message{}, "geo_information") {
		if err := db.Migrator().DropColumn(&Message{}, "geo_information"); err != nil {
			log.WithError(err).Error("failed to drop the column geo_information from Message")
		}
	}

	return nil
}
