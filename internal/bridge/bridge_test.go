package bridge

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type BridgeTestSuite struct {
	suite.Suite
}

func (s *BridgeTestSuite) SetupTest() {
	log.Logger.SetOutput(io.Discard)
}

func (s *BridgeTestSuite) TestSend() {

	s.Run("when successful", func() {
		ticker := storage.Ticker{}
		bridge := MockBridge{}
		bridge.On("Send", ticker, mock.Anything).Return(nil).Once()

		bridges := Bridges{"mock": &bridge}
		err := bridges.Send(ticker, nil)
		s.NoError(err)
		s.True(bridge.AssertExpectations(s.T()))
	})

	s.Run("when failed", func() {
		ticker := storage.Ticker{}
		bridge := MockBridge{}
		bridge.On("Send", ticker, mock.Anything).Return(errors.New("failed to send message")).Once()

		bridges := Bridges{"mock": &bridge}
		_ = bridges.Send(ticker, nil)
		s.True(bridge.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestDelete() {
	s.Run("when successful", func() {
		ticker := storage.Ticker{}
		bridge := MockBridge{}
		bridge.On("Delete", ticker, mock.Anything).Return(nil).Once()

		bridges := Bridges{"mock": &bridge}
		err := bridges.Delete(ticker, nil)
		s.NoError(err)
		s.True(bridge.AssertExpectations(s.T()))
	})

	s.Run("when failed", func() {
		ticker := storage.Ticker{}
		bridge := MockBridge{}
		bridge.On("Delete", ticker, mock.Anything).Return(errors.New("failed to delete message")).Once()

		bridges := Bridges{"mock": &bridge}
		_ = bridges.Delete(ticker, nil)
		s.True(bridge.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestRegisterBridges() {
	bridges := RegisterBridges(config.Config{}, nil)
	s.Equal(2, len(bridges))
}

func TestBrigde(t *testing.T) {
	suite.Run(t, new(BridgeTestSuite))
}
