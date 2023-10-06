package storage

import (
	"os"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/api/pagination"
)

type StormStorage struct {
	db         *storm.DB
	uploadPath string
}

func OpenDB(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		log.WithError(err).Panic("failed to open database file")
	}

	return db
}

func NewStormStorage(storagePath, uploadPath string) *StormStorage {
	return &StormStorage{
		db:         OpenDB(storagePath),
		uploadPath: uploadPath,
	}
}

func (s *StormStorage) CountUser() (int, error) {
	return s.db.Count(&User{})
}

func (s *StormStorage) FindUsers() ([]User, error) {
	users := make([]User, 0)
	err := s.db.Select().Reverse().Find(&users)
	if err != nil && err.Error() != "not found" {
		return users, err
	}

	return users, nil
}

func (s *StormStorage) FindUserByID(id int) (User, error) {
	var user User

	err := s.db.One("ID", id, &user)

	return user, err
}

func (s *StormStorage) FindUsersByTicker(ticker Ticker) ([]User, error) {
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

func (s *StormStorage) FindUserByEmail(email string) (User, error) {
	var user User

	err := s.db.One("Email", email, &user)
	return user, err
}

func (s *StormStorage) SaveUser(user *User) error {
	return s.db.Save(user)
}

func (s *StormStorage) DeleteUser(user User) error {
	return s.db.DeleteStruct(&user)
}

func (s *StormStorage) AddUsersToTicker(ticker Ticker, ids []int) error {
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

func (s *StormStorage) RemoveTickerFromUser(ticker Ticker, user User) error {
	user.RemoveTicker(ticker)

	return s.SaveUser(&user)
}

func (s *StormStorage) FindTickers() ([]Ticker, error) {
	tickers := make([]Ticker, 0)

	err := s.db.Select().Reverse().Find(&tickers)
	if err != nil && err.Error() != "not found" {
		return tickers, err
	}

	return tickers, nil
}

func (s *StormStorage) FindTickersByIDs(ids []int) ([]Ticker, error) {
	tickers := make([]Ticker, 0)

	err := s.db.Select(q.In("ID", ids)).Reverse().Find(&tickers)
	if err != nil && err.Error() != "not found" {
		return tickers, err
	}

	return tickers, nil
}

func (s *StormStorage) FindTickerByDomain(domain string) (Ticker, error) {
	var ticker Ticker

	err := s.db.One("Domain", domain, &ticker)
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}

func (s *StormStorage) FindTickerByID(id int) (Ticker, error) {
	var ticker Ticker

	err := s.db.One("ID", id, &ticker)
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}

func (s *StormStorage) SaveTicker(ticker *Ticker) error {
	return s.db.Save(ticker)
}

func (s *StormStorage) DeleteTicker(ticker Ticker) error {
	return s.db.DeleteStruct(&ticker)
}

func (s *StormStorage) FindUploadsByMessage(message Message) []Upload {
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

func (s *StormStorage) FindUploadByUUID(uuid string) (Upload, error) {
	var upload Upload
	err := s.db.One("UUID", uuid, &upload)
	return upload, err
}

func (s *StormStorage) FindUploadsByIDs(ids []int) ([]Upload, error) {
	uploads := make([]Upload, 0)
	err := s.db.Select(q.In("ID", ids)).Find(&uploads)
	return uploads, err
}

func (s *StormStorage) UploadPath() string {
	return s.uploadPath
}

func (s *StormStorage) SaveUpload(upload *Upload) error {
	return s.db.Save(upload)
}

func (s *StormStorage) DeleteUpload(upload Upload) error {
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

func (s *StormStorage) DeleteUploads(uploads []Upload) {
	for _, upload := range uploads {
		err := s.DeleteUpload(upload)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"id": upload.ID, "uuid": upload.UUID}).Error("failed to delete upload")
		}
	}
}

func (s *StormStorage) DeleteUploadsByTicker(ticker Ticker) error {
	err := s.db.Select(q.Eq("TickerID", ticker.ID)).Delete(&Upload{})
	if err != nil && err.Error() == "not found" {
		return nil
	}

	return err
}

func (s *StormStorage) FindMessage(tickerID, messageID int) (Message, error) {
	var message Message
	matcher := q.And(q.Eq("ID", messageID), q.Eq("Ticker", tickerID))
	err := s.db.Select(matcher).First(&message)
	return message, err
}

func (s *StormStorage) FindMessagesByTicker(ticker Ticker) ([]Message, error) {
	messages := make([]Message, 0)

	err := s.db.Select(q.Eq("Ticker", ticker.ID)).Reverse().Find(&messages)
	if err != nil && err.Error() == "not found" {
		return messages, nil
	}
	return messages, err
}

func (s *StormStorage) FindMessagesByTickerAndPagination(ticker Ticker, pagination pagination.Pagination) ([]Message, error) {
	messages := make([]Message, 0)

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

func (s *StormStorage) SaveMessage(message *Message) error {
	return s.db.Save(message)
}

func (s *StormStorage) DeleteMessage(message Message) error {
	uploads := s.FindUploadsByMessage(message)

	if len(uploads) > 0 {
		s.DeleteUploads(uploads)
	}

	return s.db.DeleteStruct(&message)
}

func (s *StormStorage) DeleteMessages(ticker Ticker) error {
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

func (s *StormStorage) FindSetting(name string) (Setting, error) {
	var setting Setting
	err := s.db.One("Name", name, &setting)
	if err != nil && err.Error() == "not found" {
		return setting, err
	}

	return setting, nil
}

func (s *StormStorage) GetInactiveSetting() Setting {
	setting, err := s.FindSetting(SettingInactiveName)
	if err != nil {
		return DefaultInactiveSetting()
	}

	return setting
}

func (s *StormStorage) GetRefreshIntervalSetting() Setting {
	setting, err := s.FindSetting(SettingRefreshInterval)
	if err != nil {
		return DefaultRefreshIntervalSetting()
	}

	return setting
}

func (s *StormStorage) GetRefreshIntervalSettingValue() int {
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

func (s *StormStorage) SaveInactiveSetting(inactiveSettings InactiveSettings) error {
	setting, err := s.FindSetting(SettingInactiveName)
	if err != nil {
		setting = Setting{Name: SettingInactiveName, Value: inactiveSettings}
	} else {
		setting.Value = inactiveSettings
	}

	return s.db.Save(&setting)
}

func (s *StormStorage) SaveRefreshInterval(refreshInterval float64) error {
	setting, err := s.FindSetting(SettingRefreshInterval)
	if err != nil {
		setting = Setting{Name: SettingRefreshInterval, Value: refreshInterval}
	} else {
		setting.Value = refreshInterval
	}
	return s.db.Save(&setting)
}

func (s *StormStorage) DropAll() {
	_ = s.db.Drop("User")
	_ = s.db.Drop("Message")
	_ = s.db.Drop("Ticker")
	_ = s.db.Drop("Setting")
}
