package storage

import "gorm.io/gorm"

// MigrateDB migrates the database
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&Attachment{},
		&Message{},
		&Setting{},
		&Ticker{},
		&TickerInformation{},
		&TickerMastodon{},
		&TickerTelegram{},
		&Upload{},
		&User{},
	)
}
