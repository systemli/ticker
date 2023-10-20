package storage

import (
	"encoding/json"
	"os"

	"github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlStorage struct {
	db         *gorm.DB
	uploadPath string
}

func NewSqlStorage(db *gorm.DB, uploadPath string) *SqlStorage {
	return &SqlStorage{
		db:         db,
		uploadPath: uploadPath,
	}
}

func (s *SqlStorage) CountUser() (int, error) {
	var count int64
	err := s.db.Model(&User{}).Count(&count).Error

	return int(count), err
}

func (s *SqlStorage) FindUsers() ([]User, error) {
	users := make([]User, 0)
	err := s.db.Find(&users).Error

	return users, err
}

func (s *SqlStorage) FindUserByID(id int) (User, error) {
	var user User

	err := s.db.First(&user, id).Error

	return user, err
}

func (s *SqlStorage) FindUsersByIDs(ids []int) ([]User, error) {
	users := make([]User, 0)
	err := s.db.Find(&users, ids).Error

	return users, err
}

func (s *SqlStorage) FindUsersByTicker(ticker Ticker) ([]User, error) {
	users := make([]User, 0)
	err := s.db.Model(&ticker).Association("Users").Find(&users)

	return users, err
}

func (s *SqlStorage) FindUserByEmail(email string) (User, error) {
	var user User

	err := s.db.First(&user, "email = ?", email).Error

	return user, err
}

func (s *SqlStorage) SaveUser(user *User) error {
	return s.db.Save(user).Error
}

func (s *SqlStorage) DeleteUser(user User) error {
	return s.db.Delete(&user).Error
}

func (s *SqlStorage) DeleteTickerUsers(ticker *Ticker) error {
	err := s.db.Model(ticker).Association("Users").Clear()

	return err
}

func (s *SqlStorage) DeleteTickerUser(ticker *Ticker, user *User) error {
	err := s.db.Model(ticker).Association("Users").Delete(user)

	return err
}

func (s *SqlStorage) AddTickerUser(ticker *Ticker, user *User) error {
	err := s.db.Model(ticker).Association("Users").Append(user)

	return err
}

func (s *SqlStorage) FindTickers() ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	err := s.db.Preload(clause.Associations).Find(&tickers).Error

	return tickers, err
}

func (s *SqlStorage) FindTickersByIDs(ids []int) ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	err := s.db.Preload(clause.Associations).Find(&tickers, ids).Error

	return tickers, err
}

func (s *SqlStorage) FindTickerByDomain(domain string) (Ticker, error) {
	var ticker Ticker

	err := s.db.Preload(clause.Associations).First(&ticker, "domain = ?", domain).Error

	return ticker, err
}

func (s *SqlStorage) FindTickerByID(id int) (Ticker, error) {
	var ticker Ticker

	err := s.db.Preload(clause.Associations).First(&ticker, id).Error

	return ticker, err
}

func (s *SqlStorage) SaveTicker(ticker *Ticker) error {
	return s.db.Save(ticker).Error
}

func (s *SqlStorage) DeleteTicker(ticker Ticker) error {
	return s.db.Delete(&ticker).Error
}

func (s *SqlStorage) FindUploadByUUID(uuid string) (Upload, error) {
	var upload Upload

	err := s.db.First(&upload, "uuid = ?", uuid).Error

	return upload, err
}

func (s *SqlStorage) FindUploadsByIDs(ids []int) ([]Upload, error) {
	uploads := make([]Upload, 0)
	err := s.db.Find(&uploads, ids).Error

	return uploads, err
}

func (s *SqlStorage) UploadPath() string {
	return s.uploadPath
}

func (s *SqlStorage) SaveUpload(upload *Upload) error {
	return s.db.Save(upload).Error
}

func (s *SqlStorage) DeleteUpload(upload Upload) error {
	var err error

	if err = os.Remove(upload.FullPath(s.uploadPath)); err != nil {
		log.WithError(err).WithField("upload", upload).Error("failed to delete upload file")
	}

	if err = s.db.Delete(&upload).Error; err != nil {
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
	s.db.Model(&Upload{}).Where("ticker_id = ?", ticker.ID).Find(&uploads)

	for _, upload := range uploads {
		if err := s.DeleteUpload(upload); err != nil {
			return err
		}
	}

	return nil
}

func (s *SqlStorage) FindMessage(tickerID, messageID int) (Message, error) {
	var message Message

	err := s.db.Preload(clause.Associations).First(&message, "ticker_id = ? AND id = ?", tickerID, messageID).Error

	return message, err
}

func (s *SqlStorage) FindMessagesByTicker(ticker Ticker) ([]Message, error) {
	messages := make([]Message, 0)
	err := s.db.Preload(clause.Associations).Model(&Message{}).Where("ticker_id = ?", ticker.ID).Find(&messages).Error

	return messages, err
}

func (s *SqlStorage) FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination) ([]Message, error) {
	messages := make([]Message, 0)
	query := s.db.Preload(clause.Associations).Where("ticker_id = ?", ticker.ID)

	if pagination.GetBefore() > 0 {
		query = query.Where("id < ?", pagination.GetBefore())
	} else if pagination.GetAfter() > 0 {
		query = query.Where("id > ?", pagination.GetAfter())
	}

	err := query.Order("id desc").Limit(pagination.GetLimit()).Find(&messages).Error
	return messages, err
}

func (s *SqlStorage) SaveMessage(message *Message) error {
	return s.db.Save(message).Error
}

func (s *SqlStorage) DeleteMessage(message Message) error {
	return s.db.Delete(&message).Error
}

func (s *SqlStorage) DeleteMessages(ticker Ticker) error {
	err := s.db.Where("ticker_id = ?", ticker.ID).Delete(&Message{}).Error

	return err
}

func (s *SqlStorage) GetInactiveSettings() InactiveSettings {
	var setting Setting
	err := s.db.First(&setting, "name = ?", SettingInactiveName).Error
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
	err := s.db.First(&setting, "name = ?", SettingRefreshInterval).Error
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
	err := s.db.First(&setting, "name = ?", SettingInactiveName).Error
	if err != nil {
		setting = Setting{Name: SettingInactiveName}
	}

	value, err := json.Marshal(inactiveSettings)
	if err != nil {
		return err
	}
	setting.Value = string(value)

	return s.db.Save(&setting).Error
}

func (s *SqlStorage) SaveRefreshIntervalSettings(refreshInterval RefreshIntervalSettings) error {
	var setting Setting
	err := s.db.First(&setting, "name = ?", SettingRefreshInterval).Error
	if err != nil {
		setting = Setting{Name: SettingRefreshInterval}
	}

	value, err := json.Marshal(refreshInterval)
	if err != nil {
		return err
	}
	setting.Value = string(value)

	return s.db.Save(&setting).Error
}
