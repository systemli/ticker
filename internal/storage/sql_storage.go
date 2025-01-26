package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	EqualTickerID = "ticker_id = ?"
	EqualName     = "name = ?"
)

type SqlStorage struct {
	DB         *gorm.DB
	uploadPath string
}

func NewSqlStorage(db *gorm.DB, uploadPath string) *SqlStorage {
	return &SqlStorage{
		DB:         db,
		uploadPath: uploadPath,
	}
}

func (s *SqlStorage) FindUsers(filter UserFilter, opts ...func(*gorm.DB) *gorm.DB) ([]User, error) {
	users := make([]User, 0)
	db := s.prepareDb(opts...)

	if filter.Email != nil {
		db = db.Where("email LIKE ?", fmt.Sprintf("%%%s%%", *filter.Email))
	}

	if filter.IsSuperAdmin != nil {
		db = db.Where("is_super_admin = ?", *filter.IsSuperAdmin)
	}

	err := db.Order(fmt.Sprintf("%s %s", filter.OrderBy, filter.Sort)).Find(&users).Error

	return users, err
}

func (s *SqlStorage) FindUserByID(id int, opts ...func(*gorm.DB) *gorm.DB) (User, error) {
	var user User
	db := s.prepareDb(opts...)
	err := db.First(&user, id).Error

	return user, err
}

func (s *SqlStorage) FindUsersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]User, error) {
	users := make([]User, 0)
	db := s.prepareDb(opts...)
	err := db.Find(&users, ids).Error

	return users, err
}

func (s *SqlStorage) FindUsersByTicker(ticker Ticker, opts ...func(*gorm.DB) *gorm.DB) ([]User, error) {
	users := make([]User, 0)
	db := s.prepareDb(opts...)
	err := db.Model(&ticker).Association("Users").Find(&users)

	return users, err
}

func (s *SqlStorage) FindUserByEmail(email string, opts ...func(*gorm.DB) *gorm.DB) (User, error) {
	var user User
	db := s.prepareDb(opts...)
	err := db.First(&user, "email = ?", email).Error

	return user, err
}

func (s *SqlStorage) SaveUser(user *User) error {
	if user.ID == 0 {
		return s.DB.Create(user).Error
	}

	// Replace all Tickers associations
	err := s.DB.Model(user).Association("Tickers").Replace(user.Tickers)
	if err != nil {
		log.WithError(err).WithField("user_id", user.ID).Error("failed to replace user tickers")
	}

	return s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Model(user).Updates(user.AsMap()).Error
}

func (s *SqlStorage) DeleteUser(user User) error {
	return s.DB.Delete(&user).Error
}

func (s *SqlStorage) DeleteTickerUsers(ticker *Ticker) error {
	err := s.DB.Model(ticker).Association("Users").Clear()

	return err
}

func (s *SqlStorage) DeleteTickerUser(ticker *Ticker, user *User) error {
	err := s.DB.Model(ticker).Association("Users").Delete(user)

	return err
}

func (s *SqlStorage) AddTickerUser(ticker *Ticker, user *User) error {
	err := s.DB.Model(ticker).Association("Users").Append(user)

	return err
}

func (s *SqlStorage) FindTickersByUser(user User, filter TickerFilter, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	db := s.prepareDb(opts...)

	if filter.Active != nil {
		db = db.Where("active = ?", *filter.Active)
	}

	if filter.Origin != nil {
		db = db.Joins("JOIN ticker_websites ON tickers.id = ticker_websites.ticker_id").
			Where("ticker_websites.origin LIKE ?", fmt.Sprintf("%%%s%%", *filter.Origin))
	}

	if filter.Title != nil {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", *filter.Title))
	}

	db = db.Order(fmt.Sprintf("tickers.%s %s", filter.OrderBy, filter.Sort))

	var err error
	if user.IsSuperAdmin {
		err = db.Find(&tickers).Error
	} else {
		err = db.Model(&user).Association("Tickers").Find(&tickers)
	}

	return tickers, err
}

func (s *SqlStorage) FindTickerByUserAndID(user User, id int, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error) {
	db := s.prepareDb(opts...)

	var ticker Ticker
	var err error
	if user.IsSuperAdmin {
		err = db.First(&ticker, id).Error
	} else {
		err = db.Model(&user).Association("Tickers").Find(&ticker, id)
	}

	return ticker, err
}

func (s *SqlStorage) FindTickersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	db := s.prepareDb(opts...)
	err := db.Find(&tickers, ids).Error

	return tickers, err
}

func (s *SqlStorage) FindTickerByOrigin(origin string, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error) {
	var ticker Ticker
	db := s.prepareDb(opts...)

	err := db.Joins("JOIN ticker_websites ON tickers.id = ticker_websites.ticker_id").
		Where("ticker_websites.origin = ?", origin).
		First(&ticker).Error

	return ticker, err
}

func (s *SqlStorage) FindTickerByID(id int, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error) {
	var ticker Ticker
	db := s.prepareDb(opts...)

	err := db.First(&ticker, id).Error

	return ticker, err
}

func (s *SqlStorage) SaveTicker(ticker *Ticker) error {
	if ticker.ID == 0 {
		return s.DB.Create(ticker).Error
	}

	// Replace all Users associations
	err := s.DB.Model(ticker).Association("Users").Replace(ticker.Users)
	if err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to replace ticker users")
	}

	return s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Model(ticker).Updates(ticker.AsMap()).Error
}

// DeleteTicker deletes a ticker and all associated data.
func (s *SqlStorage) DeleteTicker(ticker *Ticker) error {
	if err := s.deleteTickerAssociations(ticker); err != nil {
		return err
	}

	return s.DB.Delete(&ticker).Error
}

func (s *SqlStorage) SaveTickerWebsite(ticker *Ticker, domain string) error {
	err := s.DB.Create(&TickerWebsite{
		TickerID: ticker.ID,
		Origin:   domain,
	}).Error

	if err != nil {
		return err
	}

	return s.findTickerWebsites(ticker)
}

func (s *SqlStorage) DeleteTickerWebsite(ticker *Ticker, domain string) error {
	err := s.DB.Delete(TickerWebsite{}, "ticker_id = ? AND origin = ?", ticker.ID, domain).Error
	if err != nil {
		return err
	}

	return s.findTickerWebsites(ticker)
}

func (s *SqlStorage) DeleteTickerWebsites(ticker *Ticker) error {
	if err := s.DB.Delete(TickerWebsite{}, EqualTickerID, ticker.ID).Error; err != nil {
		return err
	}

	return s.findTickerWebsites(ticker)
}

func (s *SqlStorage) findTickerWebsites(ticker *Ticker) error {
	var websites []TickerWebsite

	if err := s.DB.Model(&ticker).Association("Websites").Find(&websites); err != nil {
		return err
	}

	ticker.Websites = websites

	return nil
}

func (s *SqlStorage) ResetTicker(ticker *Ticker) error {
	if err := s.deleteTickerAssociations(ticker); err != nil {
		return err
	}

	ticker.Active = false
	ticker.Title = "Ticker"
	ticker.Description = ""
	ticker.Information = TickerInformation{}
	ticker.Location = TickerLocation{}
	ticker.Users = make([]User, 0)

	return s.DB.Save(&ticker).Error
}

func (s *SqlStorage) deleteTickerAssociations(ticker *Ticker) error {
	if err := s.DeleteTickerUsers(ticker); err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to delete ticker users")
		return err
	}

	if err := s.DeleteTickerWebsites(ticker); err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to delete ticker websites")
		return err
	}

	if err := s.DeleteUploadsByTicker(ticker); err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to delete ticker uploads")
		return err
	}

	if err := s.DeleteMessages(ticker); err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to delete ticker messages")
		return err
	}

	if err := s.DeleteIntegrations(ticker); err != nil {
		log.WithError(err).WithField("ticker_id", ticker.ID).Error("failed to delete ticker integrations")
		return err
	}

	return nil
}

func (s *SqlStorage) DeleteIntegrations(ticker *Ticker) error {
	if err := s.DeleteMastodon(ticker); err != nil {
		return err
	}

	if err := s.DeleteTelegram(ticker); err != nil {
		return err
	}

	if err := s.DeleteBluesky(ticker); err != nil {
		return err
	}

	if err := s.DeleteSignalGroup(ticker); err != nil {
		return err
	}

	return nil
}

func (s *SqlStorage) DeleteMastodon(ticker *Ticker) error {
	ticker.Mastodon = TickerMastodon{}

	return s.DB.Delete(TickerMastodon{}, EqualTickerID, ticker.ID).Error
}

func (s *SqlStorage) DeleteTelegram(ticker *Ticker) error {
	ticker.Telegram = TickerTelegram{}

	return s.DB.Delete(TickerTelegram{}, EqualTickerID, ticker.ID).Error
}

func (s *SqlStorage) DeleteBluesky(ticker *Ticker) error {
	ticker.Bluesky = TickerBluesky{}

	return s.DB.Delete(TickerBluesky{}, EqualTickerID, ticker.ID).Error
}

func (s *SqlStorage) DeleteSignalGroup(ticker *Ticker) error {
	ticker.SignalGroup = TickerSignalGroup{}

	return s.DB.Delete(TickerSignalGroup{}, EqualTickerID, ticker.ID).Error
}

func (s *SqlStorage) FindUploadByUUID(uuid string) (Upload, error) {
	var upload Upload

	err := s.DB.First(&upload, "uuid = ?", uuid).Error

	return upload, err
}

func (s *SqlStorage) FindUploadsByIDs(ids []int) ([]Upload, error) {
	uploads := make([]Upload, 0)
	err := s.DB.Find(&uploads, ids).Error

	return uploads, err
}

func (s *SqlStorage) UploadPath() string {
	return s.uploadPath
}

func (s *SqlStorage) SaveUpload(upload *Upload) error {
	return s.DB.Save(upload).Error
}

func (s *SqlStorage) DeleteUpload(upload Upload) error {
	var err error

	if err = os.Remove(upload.FullPath(s.uploadPath)); err != nil {
		log.WithError(err).WithField("upload", upload).Error("failed to delete upload file")
	}

	if err = s.DB.Delete(&upload).Error; err != nil {
		log.WithError(err).WithField("upload", upload).Error("failed to delete upload from database")
	}

	return err
}

func (s *SqlStorage) DeleteUploads(uploads []Upload) {
	for _, upload := range uploads {
		if err := s.DeleteUpload(upload); err != nil {
			log.WithError(err).WithField("upload", upload).Error("failed to delete upload")
		}
	}
}

func (s *SqlStorage) DeleteUploadsByTicker(ticker *Ticker) error {
	uploads := make([]Upload, 0)
	s.DB.Model(&Upload{}).Where(EqualTickerID, ticker.ID).Find(&uploads)

	for _, upload := range uploads {
		if err := s.DeleteUpload(upload); err != nil {
			return err
		}
	}

	return nil
}

func (s *SqlStorage) FindMessage(tickerID, messageID int, opts ...func(*gorm.DB) *gorm.DB) (Message, error) {
	var message Message
	db := s.prepareDb(opts...)

	err := db.First(&message, "ticker_id = ? AND id = ?", tickerID, messageID).Error

	return message, err
}

func (s *SqlStorage) FindMessagesByTicker(ticker Ticker, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error) {
	messages := make([]Message, 0)
	db := s.prepareDb(opts...)

	err := db.Model(&Message{}).Where(EqualTickerID, ticker.ID).Find(&messages).Error

	return messages, err
}

func (s *SqlStorage) FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error) {
	messages := make([]Message, 0)
	db := s.prepareDb(opts...)
	query := db.Where(EqualTickerID, ticker.ID)

	if pagination.GetBefore() > 0 {
		query = query.Where("id < ?", pagination.GetBefore())
	} else if pagination.GetAfter() > 0 {
		query = query.Where("id > ?", pagination.GetAfter())
	}

	err := query.Order("id desc").Limit(pagination.GetLimit()).Find(&messages).Error
	return messages, err
}

func (s *SqlStorage) SaveMessage(message *Message) error {
	if message.ID == 0 {
		return s.DB.Create(message).Error
	}

	return s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Model(message).Updates(message.AsMap()).Error
}

func (s *SqlStorage) DeleteMessage(message Message) error {
	if len(message.Attachments) > 0 {
		err := s.DB.Where("message_id = ?", message.ID).Delete(&Attachment{}).Error
		if err != nil {
			log.WithError(err).WithField("message_id", message.ID).Error("failed to delete attachments")
		}
	}

	return s.DB.Delete(&message).Error
}

func (s *SqlStorage) DeleteMessages(ticker *Ticker) error {
	var msgIds []int
	err := s.DB.Model(&Message{}).Where(EqualTickerID, ticker.ID).Pluck("id", &msgIds).Error
	if err != nil {
		return err
	}

	err = s.DB.Where("message_id IN ?", msgIds).Delete(&Attachment{}).Error
	if err != nil {
		return err
	}

	return s.DB.Where(EqualTickerID, ticker.ID).Delete(&Message{}).Error
}

func (s *SqlStorage) GetInactiveSettings() InactiveSettings {
	var setting Setting
	err := s.DB.First(&setting, EqualName, SettingInactiveName).Error
	if err != nil {
		return DefaultInactiveSettings()
	}

	var inactiveSettings InactiveSettings
	err = json.Unmarshal([]byte(setting.Value), &inactiveSettings)
	if err != nil {
		return DefaultInactiveSettings()
	}

	return inactiveSettings
}

func (s *SqlStorage) GetRefreshIntervalSettings() RefreshIntervalSettings {
	var setting Setting
	err := s.DB.First(&setting, EqualName, SettingRefreshInterval).Error
	if err != nil {
		return DefaultRefreshIntervalSettings()
	}

	var refreshIntervalSettings RefreshIntervalSettings
	err = json.Unmarshal([]byte(setting.Value), &refreshIntervalSettings)
	if err != nil {
		return DefaultRefreshIntervalSettings()
	}

	return refreshIntervalSettings
}

func (s *SqlStorage) SaveInactiveSettings(inactiveSettings InactiveSettings) error {
	var setting Setting
	err := s.DB.First(&setting, EqualName, SettingInactiveName).Error
	if err != nil {
		setting = Setting{Name: SettingInactiveName}
	}

	value, err := json.Marshal(inactiveSettings)
	if err != nil {
		return err
	}
	setting.Value = string(value)

	return s.DB.Save(&setting).Error
}

func (s *SqlStorage) SaveRefreshIntervalSettings(refreshInterval RefreshIntervalSettings) error {
	var setting Setting
	err := s.DB.First(&setting, EqualName, SettingRefreshInterval).Error
	if err != nil {
		setting = Setting{Name: SettingRefreshInterval}
	}

	value, err := json.Marshal(refreshInterval)
	if err != nil {
		return err
	}
	setting.Value = string(value)

	return s.DB.Save(&setting).Error
}

func (s *SqlStorage) prepareDb(opts ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	db := s.DB
	for _, opt := range opts {
		db = opt(db)
	}

	return db
}

// WithPreload is a helper function to preload all associations.
func WithPreload() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(clause.Associations)
	}
}

// WithAttachments is a helper function to preload the attachments association.
func WithAttachments() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Attachments")
	}
}

// WithTickers is a helper function to preload the tickers association.
func WithTickers() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Tickers")
	}
}
