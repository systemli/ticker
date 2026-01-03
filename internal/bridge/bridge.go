package bridge

import (
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
)

var log = logger.GetWithPackage("bridge")

type Bridge interface {
	Update(ticker storage.Ticker) error
	Send(ticker storage.Ticker, message *storage.Message) error
	Delete(ticker storage.Ticker, message *storage.Message) error
}

type Bridges map[string]Bridge

func RegisterBridges(config config.Config, storage storage.Storage) Bridges {
	telegram := TelegramBridge{config, storage}
	mastodon := MastodonBridge{config, storage}
	bluesky := BlueskyBridge{config, storage}
	signalGroup := SignalGroupBridge{config, storage}
	matrix := MatrixBridge{config, storage}

	return Bridges{"telegram": &telegram, "mastodon": &mastodon, "bluesky": &bluesky, "signalGroup": &signalGroup, "matrix": &matrix}
}

func (b *Bridges) Update(ticker storage.Ticker) error {
	var err error
	for name, bridge := range *b {
		err := bridge.Update(ticker)
		if err != nil {
			log.WithError(err).WithField("bridge_name", name).Error("failed to update ticker")
		}
	}

	return err
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
