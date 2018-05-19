package storage

import (
	"github.com/asdine/storm/q"

	. "git.codecoop.org/systemli/ticker/internal/model"
	. "git.codecoop.org/systemli/ticker/internal/util"
)

//
func FindByTicker(ticker Ticker, pagination *Pagination) ([]Message, error) {
	var messages []Message

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

	err := DB.Select(matcher).OrderBy("CreationDate").Limit(pagination.GetLimit()).Reverse().Find(&messages)
	if err != nil {
		if err.Error() == "not found" {
			return messages, nil
		}
		return messages, err
	}
	return messages, nil
}
