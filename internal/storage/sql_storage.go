package storage

import (
	"encoding/json"
	"os"

	"github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (s *SqlStorage) FindUsers() ([]User, error) {
	users := make([]User, 0)
	err := s.DB.Find(&users).Error

	return users, err
}

func (s *SqlStorage) FindUserByID(id int) (User, error) {
	var user User

	err := s.DB.First(&user, id).Error

	return user, err
}

func (s *SqlStorage) FindUsersByIDs(ids []int) ([]User, error) {
	users := make([]User, 0)
	err := s.DB.Find(&users, ids).Error

	return users, err
}

func (s *SqlStorage) FindUsersByTicker(ticker Ticker) ([]User, error) {
	users := make([]User, 0)
	err := s.DB.Model(&ticker).Association("Users").Find(&users)

	return users, err
}

func (s *SqlStorage) FindUserByEmail(email string) (User, error) {
	var user User

	err := s.DB.First(&user, "email = ?", email).Error

	return user, err
}

func (s *SqlStorage) SaveUser(user *User) error {
	return s.DB.Save(user).Error
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

func (s *SqlStorage) FindTickers(opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	db := s.prepareDb(opts...)
	err := db.Find(&tickers).Error

	return tickers, err
}

func (s *SqlStorage) FindTickersByIDs(ids []int, opts ...func(*gorm.DB) *gorm.DB) ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	db := s.prepareDb(opts...)
	err := db.Find(&tickers, ids).Error

	return tickers, err
}

func (s *SqlStorage) FindTickerByDomain(domain string, opts ...func(*gorm.DB) *gorm.DB) (Ticker, error) {
	var ticker Ticker
	db := s.prepareDb(opts...)

	err := db.First(&ticker, "domain = ?", domain).Error

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

	return s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(ticker).Error
}

func (s *SqlStorage) DeleteTicker(ticker Ticker) error {
	return s.DB.Delete(&ticker).Error
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

func (s *SqlStorage) DeleteUploadsByTicker(ticker Ticker) error {
	uploads := make([]Upload, 0)
	s.DB.Model(&Upload{}).Where("ticker_id = ?", ticker.ID).Find(&uploads)

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

	err := db.Model(&Message{}).Where("ticker_id = ?", ticker.ID).Find(&messages).Error

	return messages, err
}

func (s *SqlStorage) FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination, opts ...func(*gorm.DB) *gorm.DB) ([]Message, error) {
	messages := make([]Message, 0)
	db := s.prepareDb(opts...)
	query := db.Where("ticker_id = ?", ticker.ID)

	if pagination.GetBefore() > 0 {
		query = query.Where("id < ?", pagination.GetBefore())
	} else if pagination.GetAfter() > 0 {
		query = query.Where("id > ?", pagination.GetAfter())
	}

	err := query.Order("id desc").Limit(pagination.GetLimit()).Find(&messages).Error
	return messages, err
}

func (s *SqlStorage) SaveMessage(message *Message) error {
	return s.DB.Save(message).Error
}

func (s *SqlStorage) DeleteMessage(message Message) error {
	var err error
	err = s.DB.Delete(&message).Error
	if err != nil {
		return err
	}

	if len(message.Attachments) > 0 {
		err = s.DeleteAttachmentsByMessage(message)
	}

	return err
}

func (s *SqlStorage) DeleteMessages(ticker Ticker) error {
	var msgIds []int
	err := s.DB.Model(&Message{}).Where("ticker_id = ?", ticker.ID).Pluck("id", &msgIds).Error
	if err != nil {
		return err
	}

	err = s.DB.Where("message_id IN ?", msgIds).Delete(&Attachment{}).Error
	if err != nil {
		return err
	}

	return s.DB.Where("ticker_id = ?", ticker.ID).Delete(&Message{}).Error
}

func (s *SqlStorage) DeleteAttachmentsByMessage(message Message) error {
	return s.DB.Where("message_id = ?", message.ID).Delete(&Attachment{}).Error
}

func (s *SqlStorage) GetInactiveSettings() InactiveSettings {
	var setting Setting
	err := s.DB.First(&setting, "name = ?", SettingInactiveName).Error
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
	err := s.DB.First(&setting, "name = ?", SettingRefreshInterval).Error
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
	err := s.DB.First(&setting, "name = ?", SettingInactiveName).Error
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
	err := s.DB.First(&setting, "name = ?", SettingRefreshInterval).Error
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
