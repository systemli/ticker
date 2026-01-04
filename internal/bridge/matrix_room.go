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

// Send sends a message to a Matrix room
func (mb *MatrixRoomBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	// Not implemented yet
	return nil
}

// Delete deletes a message from a Matrix room
func (mb *MatrixRoomBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	// Not implemented yet
	return nil
}
