package bridge

import (
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/matrix"
	"github.com/systemli/ticker/internal/storage"
)

type MatrixRoomBridge struct {
	config  config.Config
	storage storage.Storage
}

func (mb *MatrixRoomBridge) Update(ticker storage.Ticker) error {
	if !mb.config.Matrix.Enabled() || !ticker.Matrix.Connected() {
		return nil
	}

	// Rename room if title has changed
	err := matrix.UpdateRoomName(mb.config, ticker.Matrix.RoomID, ticker.Title)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MatrixRoomBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	if !mb.config.Matrix.Enabled() || !ticker.Matrix.Connected() || !ticker.Matrix.Active {
		return nil
	}

	eventIDs, err := matrix.SendMessage(mb.config, mb.storage, ticker.Matrix.RoomID, message)
	if err != nil {
		return err
	}

	// Store all event IDs (images + text) so we can delete them later
	message.MatrixRoom.EventIDs = eventIDs
	err = mb.storage.SaveMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MatrixRoomBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	if !mb.config.Matrix.Enabled() || !ticker.Matrix.Connected() || !ticker.Matrix.Active {
		return nil
	}

	err := matrix.DeleteMessage(mb.config, ticker.Matrix.RoomID, message.MatrixRoom.EventIDs)
	if err != nil {
		return err
	}

	return nil
}
