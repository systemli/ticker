package storage

import "gorm.io/gorm"

// Stores is the aggregate of per-aggregate stores carried by handlers, the
// bridge layer, and command-line entry points. Every field is an interface so
// individual stores can be swapped in tests.
type Stores struct {
	Users    UserStore
	Tickers  TickerStore
	Messages MessageStore
	Uploads  UploadStore
	Settings SettingsStore
	UoW      UnitOfWork
}

// NewStores wires every store against the same *SqlStorage. Production code
// uses this; tests build a Stores value with per-store mocks.
func NewStores(s *SqlStorage) Stores {
	return Stores{
		Users:    s,
		Tickers:  s,
		Messages: s,
		Uploads:  s,
		Settings: s,
		UoW:      s,
	}
}

// Tx bundles the five stores all scoped to the same database transaction.
type Tx struct {
	Users    UserStore
	Tickers  TickerStore
	Messages MessageStore
	Uploads  UploadStore
	Settings SettingsStore
}

// UnitOfWork runs a function in a database transaction. If fn returns an error
// the transaction is rolled back; otherwise it is committed.
type UnitOfWork interface {
	Do(fn func(tx Tx) error) error
}

// Do implements UnitOfWork.
func (s *SqlStorage) Do(fn func(tx Tx) error) error {
	return s.DB.Transaction(func(db *gorm.DB) error {
		return fn(Tx{
			Users:    s.WithUserTx(db),
			Tickers:  s.WithTickerTx(db),
			Messages: s.WithMessageTx(db),
			Uploads:  s.WithUploadTx(db),
			Settings: s.WithSettingsTx(db),
		})
	})
}
