package storage

import (
	"os"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/pagination"
)

var log = logrus.WithField("package", "storage")

type TickerStorage interface {
	CountUser() (int, error)
	FindUsers() ([]User, error)
	FindUserByID(id int) (User, error)
	FindUserByEmail(email string) (User, error)
	FindUsersByTicker(ticker Ticker) ([]User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error
	AddUsersToTicker(ticker Ticker, ids []int) error
	RemoveTickerFromUser(ticker Ticker, user User) error
	FindTickers() ([]Ticker, error)
	FindTickersByIDs(ids []int) ([]Ticker, error)
	FindTickerByDomain(domain string) (Ticker, error)
	FindTickerByID(id int) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker Ticker) error
	FindUploadsByMessage(message Message) []Upload
	SaveUpload(upload *Upload) error
	DeleteUpload(upload Upload) error
	DeleteUploads(uploads []Upload)
	DeleteUploadsByTicker(ticker Ticker) error
	FindMessage(tickerID, messageID int) (Message, error)
	FindMessagesByTicker(ticker Ticker) ([]Message, error)
	FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination) ([]Message, error)
	SaveMessage(message *Message) error
	DeleteMessage(message Message) error
	DeleteMessages(ticker Ticker) error
	FindSetting(name string) (Setting, error)
	GetInactiveSetting() Setting
	GetRefreshIntervalSetting() Setting
	GetRefreshIntervalSettingValue() int
	SaveInactiveSetting(inactiveSettings InactiveSettings) error
	SaveRefreshInterval(refreshInterval int) error
	FindUploadByUUID(uuid string) (Upload, error)
	FindUploadsByIDs(ids []int) ([]Upload, error)
	UploadPath() string
}

type Storage struct {
	db         *storm.DB
	uploadPath string
}

func NewStorage(storagePath, uploadPath string) *Storage {
	return &Storage{
		db:         OpenDB(storagePath),
		uploadPath: uploadPath,
	}
}

func (s *Storage) CountUser() (int, error) {
	return s.db.Count(&User{})
}

func (s *Storage) FindUsers() ([]User, error) {
	users := make([]User, 0)
	err := s.db.Select().Reverse().Find(&users)
	if err != nil && err.Error() != "not found" {
		return users, err
	}

	return users, nil
}

func (s *Storage) FindUserByID(id int) (User, error) {
	var user User

	err := s.db.One("ID", id, &user)

	return user, err
}

func (s *Storage) FindUsersByTicker(ticker Ticker) ([]User, error) {
	users := make([]User, 0)
	err := s.db.Select().Each(new(User), func(record interface{}) error {
		u := record.(*User)

		for _, id := range u.Tickers {
			if id == ticker.ID {
				users = append(users, *u)
			}
		}

		return nil
	})

	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *Storage) FindUserByEmail(email string) (User, error) {
	var user User

	err := s.db.One("Email", email, &user)
	return user, err
}

func (s *Storage) SaveUser(user *User) error {
	return s.db.Save(user)
}

func (s *Storage) DeleteUser(user User) error {
	return s.db.DeleteStruct(&user)
}

func (s *Storage) AddUsersToTicker(ticker Ticker, ids []int) error {
	users := make([]User, 0)
	err := s.db.Select(q.In("ID", ids)).Find(&users)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.IsSuperAdmin {
			continue
		}
		user.AddTicker(ticker)
		err = s.SaveUser(&user)
	}

	return err
}

func (s *Storage) RemoveTickerFromUser(ticker Ticker, user User) error {
	user.RemoveTicker(ticker)

	return s.SaveUser(&user)
}

func (s *Storage) FindTickers() ([]Ticker, error) {
	tickers := make([]Ticker, 0)

	err := s.db.Select().Reverse().Find(&tickers)
	if err != nil && err.Error() != "not found" {
		return tickers, err
	}

	return tickers, nil
}

func (s *Storage) FindTickersByIDs(ids []int) ([]Ticker, error) {
	tickers := make([]Ticker, 0)

	err := s.db.Select(q.In("ID", ids)).Reverse().Find(&tickers)
	if err != nil && err.Error() != "not found" {
		return tickers, err
	}

	return tickers, nil
}

func (s *Storage) FindTickerByDomain(domain string) (Ticker, error) {
	var ticker Ticker

	err := s.db.One("Domain", domain, &ticker)
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}

func (s *Storage) FindTickerByID(id int) (Ticker, error) {
	var ticker Ticker

	err := s.db.One("ID", id, &ticker)
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}

func (s *Storage) SaveTicker(ticker *Ticker) error {
	return s.db.Save(ticker)
}

func (s *Storage) DeleteTicker(ticker Ticker) error {
	return s.db.DeleteStruct(&ticker)
}

func (s *Storage) FindUploadsByMessage(message Message) []Upload {
	uploads := make([]Upload, 0)

	if len(message.Attachments) > 0 {
		var uuids []string
		for _, attachment := range message.Attachments {
			uuids = append(uuids, attachment.UUID)
		}
		err := s.db.Select(q.In("UUID", uuids)).Find(&uploads)
		if err != nil {
			log.WithError(err).Error("failed to find uploads for message")
		}
	}

	return uploads
}

func (s *Storage) FindUploadByUUID(uuid string) (Upload, error) {
	var upload Upload
	err := s.db.One("UUID", uuid, &upload)
	return upload, err
}

func (s *Storage) FindUploadsByIDs(ids []int) ([]Upload, error) {
	uploads := make([]Upload, 0)
	err := s.db.Select(q.In("ID", ids)).Find(&uploads)
	return uploads, err
}

func (s *Storage) UploadPath() string {
	return s.uploadPath
}

func (s *Storage) SaveUpload(upload *Upload) error {
	return s.db.Save(upload)
}

func (s *Storage) DeleteUpload(upload Upload) error {
	var err error
	//TODO: Rework with afero.FS from Config
	if err = os.Remove(upload.FullPath(s.uploadPath)); err != nil {
		log.WithError(err).WithField("upload", upload).Error("failed to delete upload file")
	}

	if err = s.db.DeleteStruct(&upload); err != nil {
		log.WithError(err).WithField("upload", upload).Error("failed to delete upload")
	}

	return err
}

func (s *Storage) DeleteUploads(uploads []Upload) {
	for _, upload := range uploads {
		err := s.DeleteUpload(upload)
		log.WithError(err).WithFields(logrus.Fields{"id": upload.ID, "uuid": upload.UUID}).Error("failed to delete upload")
	}
}

func (s *Storage) DeleteUploadsByTicker(ticker Ticker) error {
	err := s.db.Select(q.Eq("TickerID", ticker.ID)).Delete(&Upload{})
	if err != nil && err.Error() == "not found" {
		return nil
	}

	return err
}

func (s *Storage) FindMessage(tickerID, messageID int) (Message, error) {
	var message Message
	matcher := q.And(q.Eq("ID", messageID), q.Eq("Ticker", tickerID))
	err := s.db.Select(matcher).First(&message)
	return message, err
}

func (s *Storage) FindMessagesByTicker(ticker Ticker) ([]Message, error) {
	messages := make([]Message, 0)

	err := s.db.Select(q.Eq("Ticker", ticker.ID)).Reverse().Find(&messages)
	if err != nil && err.Error() == "not found" {
		return messages, nil
	}
	return messages, err
}

func (s *Storage) FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination) ([]Message, error) {
	messages := make([]Message, 0)

	if !ticker.Active {
		return messages, nil
	}

	matcher := q.Eq("Ticker", ticker.ID)
	if pagination.GetBefore() != 0 {
		matcher = q.And(q.Eq("Ticker", ticker.ID), q.Lt("ID", pagination.GetBefore()))
	}
	if pagination.GetAfter() != 0 {
		matcher = q.And(q.Eq("Ticker", ticker.ID), q.Gt("ID", pagination.GetAfter()))
	}

	err := s.db.Select(matcher).OrderBy("CreationDate").Limit(pagination.GetLimit()).Reverse().Find(&messages)
	if err != nil && err.Error() == "not found" {
		return messages, nil
	}
	return messages, err
}

func (s *Storage) SaveMessage(message *Message) error {
	return s.db.Save(message)
}

func (s *Storage) DeleteMessage(message Message) error {
	uploads := s.FindUploadsByMessage(message)

	if len(uploads) > 0 {
		s.DeleteUploads(uploads)
	}

	return s.db.DeleteStruct(&message)
}

func (s *Storage) DeleteMessages(ticker Ticker) error {
	var messages []Message
	if err := s.db.Find("Ticker", ticker.ID, &messages); err != nil {
		log.WithField("error", err).WithField("ticker", ticker.ID).Error("failed find messages for ticker")
		return err
	}

	for _, message := range messages {
		_ = s.DeleteMessage(message)
	}

	return nil
}

func (s *Storage) FindSetting(name string) (Setting, error) {
	var setting Setting
	err := s.db.One("Name", name, &setting)
	if err != nil && err.Error() == "not found" {
		return setting, err
	}

	return setting, nil
}

func (s *Storage) GetInactiveSetting() Setting {
	setting, err := s.FindSetting(SettingInactiveName)
	if err != nil {
		return DefaultInactiveSetting()
	}

	return setting
}

func (s *Storage) GetRefreshIntervalSetting() Setting {
	setting, err := s.FindSetting(SettingRefreshInterval)
	if err != nil {
		return DefaultRefreshIntervalSetting()
	}

	return setting
}

func (s *Storage) GetRefreshIntervalSettingValue() int {
	setting := s.GetRefreshIntervalSetting()

	var value int
	switch sv := setting.Value.(type) {
	case float64:
		value = int(sv)
	default:
		value = sv.(int)
	}

	return value
}

func (s *Storage) SaveInactiveSetting(inactiveSettings InactiveSettings) error {
	setting, err := s.FindSetting(SettingInactiveName)
	if err != nil {
		setting = Setting{Name: SettingInactiveName, Value: inactiveSettings}
	} else {
		setting.Value = inactiveSettings
	}

	return s.db.Save(&setting)
}

func (s *Storage) SaveRefreshInterval(refreshInterval int) error {
	setting, err := s.FindSetting(SettingRefreshInterval)
	if err != nil {
		setting = Setting{Name: SettingRefreshInterval, Value: refreshInterval}
	} else {
		setting.Value = refreshInterval
	}
	return s.db.Save(&setting)
}

func (s *Storage) DropAll() {
	_ = s.db.Drop("User")
	_ = s.db.Drop("Message")
	_ = s.db.Drop("Ticker")
	_ = s.db.Drop("Setting")
}
