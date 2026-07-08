package storage

import (
	"github.com/systemli/ticker/internal/api/pagination"
)

// MessageStore covers Message lookups, paginated reads, and writes.
type MessageStore interface {
	FindMessage(tickerID, messageID int, opts ...QueryOpt) (Message, error)
	FindMessagesByTicker(ticker Ticker, opts ...QueryOpt) ([]Message, error)
	FindMessagesByTickerAndPagination(ticker Ticker, p pagination.Pagination, opts ...QueryOpt) ([]Message, error)
	SaveMessage(message *Message) error
	DeleteMessage(message Message) error
	DeleteMessages(ticker *Ticker) error
}
