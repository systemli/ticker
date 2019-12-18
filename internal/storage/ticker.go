package storage

import (
	. "github.com/systemli/ticker/internal/model"
)

//FindTicker returns a Ticker for a given domain.
func FindTicker(domain string) (*Ticker, error) {
	var ticker Ticker

	err := DB.One("Domain", domain, &ticker)
	if err != nil {
		return &ticker, err
	}

	return &ticker, nil
}

//GetTicker returns a Ticker for given id.
func GetTicker(id int) (*Ticker, error) {
	var ticker Ticker

	err := DB.One("ID", id, &ticker)
	if err != nil {
		return &ticker, err
	}

	return &ticker, nil
}
