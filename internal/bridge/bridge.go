package bridge

import (
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var log = logrus.WithField("package", "bridge")

type Bridge interface {
	Send(ticker storage.Ticker, message *storage.Message) error
	Delete(ticker storage.Ticker, message *storage.Message) error
}

type Bridges map[string]Bridge

func RegisterBridges(config config.Config, storage storage.Storage) Bridges {
	telegram := TelegramBridge{config, storage}
	mastodon := MastodonBridge{config, storage}

	return Bridges{"telegram": &telegram, "mastodon": &mastodon}
}

func (b *Bridges) Send(ticker storage.Ticker, message *storage.Message) error {
	var err error
	for name, bridge := range *b {
		err := bridge.Send(ticker, message)
		if err != nil {
			log.WithError(err).WithField("bridge_name", name).Error("failed to send message")
		}
	}

	return err
}

func (b *Bridges) Delete(ticker storage.Ticker, message *storage.Message) error {
	var err error
	for name, bridge := range *b {
		err := bridge.Delete(ticker, message)
		if err != nil {
			log.WithError(err).WithField("bridge_name", name).Error("failed to delete message")
		}
	}

	return err
}
