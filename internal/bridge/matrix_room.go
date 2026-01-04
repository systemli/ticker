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

	eventID, err := matrix.SendMessage(mb.config, ticker.Matrix.RoomID, message.Text)
	if err != nil {
		return err
	}

	// Store the event ID so we can delete the message later
	message.MatrixRoom.EventID = eventID
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

	err := matrix.DeleteMessage(mb.config, ticker.Matrix.RoomID, message.MatrixRoom.EventID)
	if err != nil {
		return err
	}

	return nil
}
