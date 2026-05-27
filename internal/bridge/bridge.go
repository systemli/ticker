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

func RegisterBridges(config config.Config, stores storage.Stores) Bridges {
	telegram := TelegramBridge{config: config, uploads: stores.Uploads, settings: stores.Settings}
	mastodon := MastodonBridge{config: config, uploads: stores.Uploads}
	bluesky := BlueskyBridge{config: config, uploads: stores.Uploads}
	signalGroup := SignalGroupBridge{config: config, uploads: stores.Uploads, settings: stores.Settings}

	return Bridges{"telegram": &telegram, "mastodon": &mastodon, "bluesky": &bluesky, "signalGroup": &signalGroup}
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
