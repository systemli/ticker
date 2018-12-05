package storage

import (
	. "github.com/systemli/ticker/internal/model"
)

//Find Ticker Configuration by domain
func FindTicker(domain string) (*Ticker, error) {
	var ticker Ticker

	err := DB.One("Domain", domain, &ticker)
	if err != nil {
		return &ticker, err
	}

	return &ticker, nil
}
