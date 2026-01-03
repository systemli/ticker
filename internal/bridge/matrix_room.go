package bridge

import (
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type MatrixBridge struct {
	config  config.Config
	storage storage.Storage
}

// Update is called when a ticker is updated
func (mb *MatrixBridge) Update(ticker storage.Ticker) error {
	if !mb.config.Matrix.Enabled() || !ticker.Matrix.Connected() {
		return nil
	}

	// Call

	return nil
}

// Send sends a message to a Matrix room
func (mb *MatrixBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	// Not implemented yet
	return nil
}

// Delete deletes a message from a Matrix room
func (mb *MatrixBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	// Not implemented yet
	return nil
}
