package bridge

import (
	"errors"
	"io"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var tickerWithoutBridges storage.Ticker
var tickerWithBridges storage.Ticker

var messageWithoutBridges storage.Message
var messageWithBridges storage.Message

type BridgeTestSuite struct {
	suite.Suite
}

func (s *BridgeTestSuite) SetupTest() {
	log.Logger.SetOutput(io.Discard)
	gock.DisableNetworking()
	defer gock.Off()

	tickerWithoutBridges = storage.Ticker{
		Mastodon: storage.TickerMastodon{
			Active: false,
		},
		Telegram: storage.TickerTelegram{
			Active: false,
		},
	}
	tickerWithBridges = storage.Ticker{
		Mastodon: storage.TickerMastodon{
			Active:      true,
			Server:      "https://systemli.social",
			Token:       "token",
			Secret:      "secret",
			AccessToken: "access_token",
		},
		Telegram: storage.TickerTelegram{
			Active:      true,
			ChannelName: "channel",
		},
		Bluesky: storage.TickerBluesky{
			Active: true,
			Handle: "handle",
			AppKey: "app_key",
		},
		SignalGroup: storage.TickerSignalGroup{
			Active:           true,
			GroupID:          "group_id",
			GroupName:        "group_name",
			GroupDescription: "group_description",
		},
	}
	messageWithoutBridges = storage.Message{
		Text: "Hello World",
		Attachments: []storage.Attachment{
			{
				UUID: "123",
			},
		},
	}
	messageWithBridges = storage.Message{
		Text: "Hello World",
		Attachments: []storage.Attachment{
			{
				UUID: "123",
			},
			{
				UUID: "456",
			},
		},
		Mastodon: storage.MastodonMeta{
			ID: "123",
		},
		Telegram: storage.TelegramMeta{
			Messages: []tgbotapi.Message{
				{
					MessageID: 123,
					Chat: &tgbotapi.Chat{
						ID: 123,
					},
				},
			},
		},
		Bluesky: storage.BlueskyMeta{
			Handle: "handle",
			Uri:    "at://did:plc:sample-uri",
			Cid:    "cid",
		},
		SignalGroup: storage.SignalGroupMeta{
			Timestamp: 123,
		},
	}
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
	s.Equal(4, len(bridges))
}

func TestBrigde(t *testing.T) {
	suite.Run(t, new(BridgeTestSuite))
}
