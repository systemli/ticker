package storage

import "gorm.io/gorm"

// MigrateDB migrates the database
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&Ticker{},
		&TickerMastodon{},
		&TickerTelegram{},
		&TickerBluesky{},
		&TickerSignalGroup{},
		&User{},
		&Setting{},
		&Upload{},
		&Message{},
		&Attachment{},
	)
}
