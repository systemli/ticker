package storage

import (
	"github.com/asdine/storm/q"
	log "github.com/sirupsen/logrus"

	"github.com/systemli/ticker/internal/bridge"
	. "github.com/systemli/ticker/internal/model"
	. "github.com/systemli/ticker/internal/util"
)

//
func FindByTicker(ticker *Ticker, pagination *Pagination) ([]Message, error) {
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

//DeleteMessage removes a Message for a Ticker
func DeleteMessage(ticker *Ticker, message *Message) error {
	uploads, err := FindUploadsByMessage(message)
	if err != nil {
		log.WithField("error", err).WithField("message", message).Error("failed to find uploads")
		return err
	}

	DeleteUploads(uploads)

	if message.Tweet.ID != "" {
		err := bridge.Twitter.Delete(*ticker, message.Tweet.ID)
		if err != nil {
			log.WithField("error", err).WithField("message", message).Error("failed to delete tweet")
		}
	}

	err = DB.DeleteStruct(message)
	if err != nil {
		log.WithField("error", err).WithField("message", message).Error("failed to delete message")
		return err
	}

	return nil
}

//DeleteMessages removes all messages for a Ticker.
func DeleteMessages(ticker *Ticker) error {
	var messages []*Message
	err := DB.Find("ID", ticker.ID, &messages)
	if err != nil {
		log.WithField("error", err).WithField("ticker", ticker).Error("failed find messages for ticker")
		return err
	}

	for _, message := range messages {
		_ = DeleteMessage(ticker, message)
	}

	return nil
}
