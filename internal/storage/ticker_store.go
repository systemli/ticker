package storage

import "fmt"

// TickerStore covers Ticker CRUD, its sub-collections (websites), the integration
// configs (Telegram/Mastodon/Bluesky/SignalGroup), and User<->Ticker membership.
type TickerStore interface {
	// Tickers
	FindTickersByUser(user User, filter TickerFilter, opts ...QueryOpt) ([]Ticker, error)
	FindTickerByUserAndID(user User, id int, opts ...QueryOpt) (Ticker, error)
	FindTickersByIDs(ids []int, opts ...QueryOpt) ([]Ticker, error)
	FindTickerByOrigin(origin string, opts ...QueryOpt) (Ticker, error)
	FindTickerByID(id int, opts ...QueryOpt) (Ticker, error)
	SaveTicker(ticker *Ticker) error
	DeleteTicker(ticker *Ticker) error
	ResetTicker(ticker *Ticker) error

	// Websites (sub-collection)
	SaveTickerWebsites(ticker *Ticker, websites []TickerWebsite) error
	DeleteTickerWebsites(ticker *Ticker) error

	// Integration configs
	ClearIntegration(ticker *Ticker, integration Integration) error
	ClearIntegrations(ticker *Ticker) error

	// Membership (M:N with User)
	FindUsersByTicker(ticker Ticker, opts ...QueryOpt) ([]User, error)
	AddTickerUser(ticker *Ticker, user *User) error
	DeleteTickerUser(ticker *Ticker, user *User) error
	DeleteTickerUsers(ticker *Ticker) error
}

// integrationClearers maps each Integration to the function that resets its
// field on the ticker and deletes its configuration row. Table-driven so an
// unmatched Integration is a lookup miss (returns an error) rather than a
// silent no-op, and adding an integration means adding one entry here.
var integrationClearers = map[Integration]func(s *SqlStorage, ticker *Ticker) error{
	IntegrationTelegram: func(s *SqlStorage, ticker *Ticker) error {
		ticker.Telegram = TickerTelegram{}
		return s.DB.Delete(TickerTelegram{}, EqualTickerID, ticker.ID).Error
	},
	IntegrationMastodon: func(s *SqlStorage, ticker *Ticker) error {
		ticker.Mastodon = TickerMastodon{}
		return s.DB.Delete(TickerMastodon{}, EqualTickerID, ticker.ID).Error
	},
	IntegrationBluesky: func(s *SqlStorage, ticker *Ticker) error {
		ticker.Bluesky = TickerBluesky{}
		return s.DB.Delete(TickerBluesky{}, EqualTickerID, ticker.ID).Error
	},
	IntegrationSignalGroup: func(s *SqlStorage, ticker *Ticker) error {
		ticker.SignalGroup = TickerSignalGroup{}
		return s.DB.Delete(TickerSignalGroup{}, EqualTickerID, ticker.ID).Error
	},
}

// ClearIntegration removes the configuration row for a single integration on
// the given ticker.
func (s *SqlStorage) ClearIntegration(ticker *Ticker, integration Integration) error {
	clear, ok := integrationClearers[integration]
	if !ok {
		return fmt.Errorf("unknown integration: %q", integration)
	}
	return clear(s, ticker)
}

// ClearIntegrations clears every configured integration on the ticker.
func (s *SqlStorage) ClearIntegrations(ticker *Ticker) error {
	for _, integration := range AllIntegrations {
		if err := s.ClearIntegration(ticker, integration); err != nil {
			return err
		}
	}
	return nil
}

