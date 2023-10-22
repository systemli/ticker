package legacy

import (
	"time"

	"github.com/systemli/ticker/internal/storage"
)

type Migration struct {
	oldStorage *LegacyStorage
	newStorage *storage.SqlStorage
}

func NewMigration(oldStorage *LegacyStorage, newStorage *storage.SqlStorage) *Migration {
	return &Migration{
		oldStorage: oldStorage,
		newStorage: newStorage,
	}
}

func (m *Migration) Do() error {
	tickers, err := m.oldStorage.FindTickers()
	if err != nil {
		log.WithError(err).Error("Unable to find tickers")
	}

	for _, oldTicker := range tickers {
		ticker := storage.Ticker{
			ID:          oldTicker.ID,
			CreatedAt:   oldTicker.CreationDate,
			UpdatedAt:   time.Now(),
			Domain:      oldTicker.Domain,
			Title:       oldTicker.Title,
			Description: oldTicker.Description,
			Active:      oldTicker.Active,
			Information: storage.TickerInformation{
				Author:   oldTicker.Information.Author,
				URL:      oldTicker.Information.URL,
				Email:    oldTicker.Information.Email,
				Twitter:  oldTicker.Information.Twitter,
				Facebook: oldTicker.Information.Facebook,
				Telegram: oldTicker.Information.Telegram,
			},
			Telegram: storage.TickerTelegram{
				CreatedAt:   oldTicker.CreationDate,
				UpdatedAt:   time.Now(),
				Active:      oldTicker.Telegram.Active,
				ChannelName: oldTicker.Telegram.ChannelName,
			},
			Mastodon: storage.TickerMastodon{
				CreatedAt:   oldTicker.CreationDate,
				UpdatedAt:   time.Now(),
				Active:      oldTicker.Mastodon.Active,
				Server:      oldTicker.Mastodon.Server,
				Token:       oldTicker.Mastodon.Token,
				Secret:      oldTicker.Mastodon.Secret,
				AccessToken: oldTicker.Mastodon.AccessToken,
				User: storage.MastodonUser{
					Username:    oldTicker.Mastodon.User.Username,
					Avatar:      oldTicker.Mastodon.User.Avatar,
					DisplayName: oldTicker.Mastodon.User.DisplayName,
				},
			},
			Location: storage.TickerLocation{
				Lat: oldTicker.Location.Lat,
				Lon: oldTicker.Location.Lon,
			},
		}

		if err := m.newStorage.DB.Create(&ticker).Error; err != nil {
			log.WithError(err).WithField("ticker_id", ticker.ID).Error("Unable to save ticker")
		}

		messages, err := m.oldStorage.FindMessageByTickerID(oldTicker.ID)
		if err != nil {
			log.WithError(err).WithField("ticker_id", oldTicker.ID).Error("Unable to find messages for ticker")
			continue
		}

		for _, oldMessage := range messages {
			attachments := make([]storage.Attachment, 0)
			for _, oldAttachment := range oldMessage.Attachments {
				attachment := storage.Attachment{
					CreatedAt:   oldMessage.CreationDate,
					UpdatedAt:   oldMessage.CreationDate,
					MessageID:   oldMessage.ID,
					UUID:        oldAttachment.UUID,
					Extension:   oldAttachment.Extension,
					ContentType: oldAttachment.ContentType,
				}
				attachments = append(attachments, attachment)
			}

			message := storage.Message{
				ID:             oldMessage.ID,
				CreatedAt:      oldMessage.CreationDate,
				UpdatedAt:      time.Now(),
				TickerID:       oldMessage.Ticker,
				Text:           oldMessage.Text,
				Attachments:    attachments,
				GeoInformation: oldMessage.GeoInformation,
				Telegram: storage.TelegramMeta{
					Messages: oldMessage.Telegram.Messages,
				},
				Mastodon: storage.MastodonMeta{
					ID:  string(oldMessage.Mastodon.ID),
					URI: oldMessage.Mastodon.URI,
					URL: oldMessage.Mastodon.URL,
				},
			}

			if err := m.newStorage.DB.Create(&message).Error; err != nil {
				log.WithError(err).WithField("message_id", message.ID).Error("Unable to save message")
				continue
			}
		}

	}

	users, err := m.oldStorage.FindUsers()
	if err != nil {
		log.WithError(err).Error("Unable to find users")
	}

	for _, oldUser := range users {
		user := storage.User{
			ID:                oldUser.ID,
			CreatedAt:         oldUser.CreationDate,
			UpdatedAt:         time.Now(),
			Email:             oldUser.Email,
			EncryptedPassword: oldUser.EncryptedPassword,
			IsSuperAdmin:      oldUser.IsSuperAdmin,
		}

		if err := m.newStorage.DB.Create(&user).Error; err != nil {
			log.WithError(err).WithField("user_id", user.ID).Error("Unable to save user")
			continue
		}

		for _, tickerID := range oldUser.Tickers {
			ticker, err := m.newStorage.FindTickerByID(tickerID)
			if err != nil {
				log.WithError(err).WithField("ticker_id", tickerID).Warn("Unable to find ticker")
				continue
			}
			err = m.newStorage.AddTickerUser(&ticker, &user)
			if err != nil {
				log.WithError(err).WithField("ticker_id", tickerID).WithField("user_id", user.ID).Error("Unable to add ticker to user")
				continue
			}
		}
	}

	uploads, err := m.oldStorage.FindUploads()
	if err != nil {
		log.WithError(err).Error("Unable to find uploads")
	}

	for _, oldUpload := range uploads {
		upload := storage.Upload{
			ID:          oldUpload.ID,
			CreatedAt:   oldUpload.CreationDate,
			UpdatedAt:   time.Now(),
			TickerID:    oldUpload.TickerID,
			UUID:        oldUpload.UUID,
			Extension:   oldUpload.Extension,
			ContentType: oldUpload.ContentType,
		}

		if err := m.newStorage.DB.Create(&upload).Error; err != nil {
			log.WithError(err).WithField("upload_id", upload.ID).Error("Unable to save upload")
			continue
		}
	}

	refreshIntervalSetting, err := m.oldStorage.FindSetting("refresh_interval")
	if err != nil {
		log.WithError(err).Warn("Unable to find refresh_interval setting")
	} else {
		var value int
		switch refreshIntervalSetting.Value.(type) {
		case float64:
			value = int(refreshIntervalSetting.Value.(float64))
		case int:
			value = refreshIntervalSetting.Value.(int)
		}
		err = m.newStorage.SaveRefreshIntervalSettings(storage.RefreshIntervalSettings{
			RefreshInterval: value,
		})
		if err != nil {
			log.WithError(err).Error("Unable to save refresh_interval setting")
		}
	}

	inactiveSettings, err := m.oldStorage.FindSetting("inactive_settings")
	if err != nil {
		log.WithError(err).Warn("Unable to find inactive_settings setting")
	} else {
		err = m.newStorage.SaveInactiveSettings(storage.InactiveSettings{
			Headline:    inactiveSettings.Value.(map[string]interface{})["headline"].(string),
			SubHeadline: inactiveSettings.Value.(map[string]interface{})["sub_headline"].(string),
			Description: inactiveSettings.Value.(map[string]interface{})["description"].(string),
			Author:      inactiveSettings.Value.(map[string]interface{})["author"].(string),
			Email:       inactiveSettings.Value.(map[string]interface{})["email"].(string),
			Twitter:     inactiveSettings.Value.(map[string]interface{})["twitter"].(string),
		})
		if err != nil {
			log.WithError(err).Error("Unable to save inactive_settings setting")
		}
	}

	return nil
}
