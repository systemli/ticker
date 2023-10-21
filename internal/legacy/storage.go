package legacy

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type LegacyStorage struct {
	db *storm.DB
}

func NewLegacyStorage(db *storm.DB) *LegacyStorage {
	return &LegacyStorage{db: db}
}

func (s *LegacyStorage) FindTickers() ([]Ticker, error) {
	tickers := make([]Ticker, 0)
	err := s.db.Select().Reverse().Find(&tickers)
	if err != nil && err.Error() != "not found" {
		return tickers, err
	}

	return tickers, nil
}

func (s *LegacyStorage) FindUsers() ([]User, error) {
	users := make([]User, 0)
	err := s.db.Select().Reverse().Find(&users)
	if err != nil && err.Error() != "not found" {
		return users, err
	}

	return users, nil
}

func (s *LegacyStorage) FindUploads() ([]Upload, error) {
	uploads := make([]Upload, 0)
	err := s.db.Select().Reverse().Find(&uploads)
	if err != nil && err.Error() != "not found" {
		return uploads, err
	}

	return uploads, nil
}

func (s *LegacyStorage) FindMessageByTickerID(id int) ([]Message, error) {
	messages := make([]Message, 0)
	err := s.db.Select(q.Eq("Ticker", id)).Reverse().Find(&messages)
	if err != nil && err.Error() == "not found" {
		return messages, nil
	}

	return messages, err
}
