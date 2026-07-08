package storage

// Stores is the aggregate of per-aggregate stores carried by handlers, the
// bridge layer, and command-line entry points. Every field is an interface so
// individual stores can be swapped in tests.
type Stores struct {
	Users    UserStore
	Tickers  TickerStore
	Messages MessageStore
	Uploads  UploadStore
	Settings SettingsStore
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
	}
}
