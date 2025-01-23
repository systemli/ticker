package storage

import (
	"fmt"
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

	// Migrate all Ticker.Origin to TickerWebsite
	var tickers []Ticker
	if err := db.Find(&tickers).Error; err != nil {
		return err
	}

	for _, ticker := range tickers {
		if ticker.Domain != "" {
			if err := db.Create(&TickerWebsite{
				TickerID: ticker.ID,
				Origin:   migrateDomain(ticker.Domain),
			}).Error; err != nil {
				return err
			}

			if err := db.Model(&ticker).Update("Domain", "").Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func migrateDomain(old string) string {
	if old == "localhost" {
		return "http://localhost"
	}

	return fmt.Sprintf("https://%s", old)
}
