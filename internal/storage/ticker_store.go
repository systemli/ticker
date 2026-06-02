package storage

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

// ClearIntegration removes the configuration row for a single integration on
// the given ticker.
func (s *SqlStorage) ClearIntegration(ticker *Ticker, integration Integration) error {
	switch integration {
	case IntegrationTelegram:
		ticker.Telegram = TickerTelegram{}
		return s.DB.Delete(TickerTelegram{}, EqualTickerID, ticker.ID).Error
	case IntegrationMastodon:
		ticker.Mastodon = TickerMastodon{}
		return s.DB.Delete(TickerMastodon{}, EqualTickerID, ticker.ID).Error
	case IntegrationBluesky:
		ticker.Bluesky = TickerBluesky{}
		return s.DB.Delete(TickerBluesky{}, EqualTickerID, ticker.ID).Error
	case IntegrationSignalGroup:
		ticker.SignalGroup = TickerSignalGroup{}
		return s.DB.Delete(TickerSignalGroup{}, EqualTickerID, ticker.ID).Error
	}
	return nil
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

