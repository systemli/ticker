package bridge

import (
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/signal"
	"github.com/systemli/ticker/internal/storage"
)

type SignalGroupBridge struct {
	config  config.Config
	storage storage.Storage
}

func (sb *SignalGroupBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	if !sb.config.SignalGroup.Enabled() || !ticker.SignalGroup.Connected() || !ticker.SignalGroup.Active {
		return nil
	}

	err := signal.SendGroupMessage(sb.config, sb.storage, ticker.SignalGroup.GroupID, message)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SignalGroupBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	if !sb.config.SignalGroup.Enabled() || !ticker.SignalGroup.Connected() || !ticker.SignalGroup.Active || message.SignalGroup.Timestamp == nil {
		return nil
	}

	err := signal.DeleteMessage(sb.config, ticker.SignalGroup.GroupID, message)
	if err != nil {
		return err
	}

	return nil

}
