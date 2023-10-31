package bridge

import (
	"errors"
	"io"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

var tickerWithoutMastodon storage.Ticker
var tickerWithMastodon storage.Ticker

var messageWithoutMastodon storage.Message
var messageWithMastodon storage.Message

type MastodonBridgeTestSuite struct {
	suite.Suite
}

func (s *MastodonBridgeTestSuite) SetupTest() {
	log.Logger.SetOutput(io.Discard)
	gock.DisableNetworking()
	defer gock.Off()

	tickerWithoutMastodon = storage.Ticker{
		Mastodon: storage.TickerMastodon{
			Active: false,
		},
	}
	tickerWithMastodon = storage.Ticker{
		Mastodon: storage.TickerMastodon{
			Active:      true,
			Server:      "https://systemli.social",
			Token:       "token",
			Secret:      "secret",
			AccessToken: "access_token",
		},
	}
	messageWithoutMastodon = storage.Message{
		Text: "Hello World",
		Attachments: []storage.Attachment{
			{
				UUID: "123",
			},
		},
	}
	messageWithMastodon = storage.Message{
		Text: "Hello World",
		Attachments: []storage.Attachment{
			{
				UUID: "123",
			},
		},
		Mastodon: storage.MastodonMeta{
			ID: "123",
		},
	}
}

func (s *MastodonBridgeTestSuite) TestSend() {
	s.Run("when mastodon is inactive", func() {
		bridge := s.bridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Send(tickerWithoutMastodon, &messageWithoutMastodon)
		s.NoError(err)
	})

	s.Run("when mastodon is active but upload cant not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.bridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(200).
			JSON(map[string]string{
				"id":  "123",
				"uri": "https://systemli.social/@systemli/123",
				"url": "https://systemli.social/@systemli/123",
			})

		err := bridge.Send(tickerWithMastodon, &messageWithoutMastodon)
		s.NoError(err)
		s.Equal("123", messageWithoutMastodon.Mastodon.ID)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutMastodon.Mastodon.URI)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutMastodon.Mastodon.URL)
		s.True(gock.IsDone())
	})

	s.Run("when mastodon is active and upload is not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, errors.New("upload not found")).Once()
		bridge := s.bridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(200).
			JSON(map[string]string{
				"id":  "123",
				"uri": "https://systemli.social/@systemli/123",
				"url": "https://systemli.social/@systemli/123",
			})

		err := bridge.Send(tickerWithMastodon, &messageWithoutMastodon)
		s.NoError(err)
		s.Equal("123", messageWithoutMastodon.Mastodon.ID)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutMastodon.Mastodon.URI)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutMastodon.Mastodon.URL)
		s.True(gock.IsDone())
		s.True(gock.IsDone())
	})

	s.Run("when mastodon is active but post status fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.bridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(500)

		err := bridge.Send(tickerWithMastodon, &messageWithoutMastodon)
		s.Error(err)
		s.True(gock.IsDone())
	})
}

func (s *MastodonBridgeTestSuite) TestDelete() {
	s.Run("when message has no mastodon meta", func() {
		bridge := s.bridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithMastodon, &messageWithoutMastodon)
		s.NoError(err)
	})

	s.Run("when mastodon is inactive", func() {
		bridge := s.bridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithoutMastodon, &messageWithMastodon)
		s.Error(err)
		s.Equal("unable to delete the status", err.Error())
	})

	s.Run("when mastodon is active", func() {
		bridge := s.bridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://systemli.social").
			Delete("/api/v1/statuses/123").
			Reply(200)

		err := bridge.Delete(tickerWithMastodon, &messageWithMastodon)
		s.NoError(err)
		s.True(gock.IsDone())
	})
}

func (s *MastodonBridgeTestSuite) bridge(config config.Config, storage storage.Storage) *MastodonBridge {
	return &MastodonBridge{
		config:  config,
		storage: storage,
	}
}

func TestMastodonBridge(t *testing.T) {
	suite.Run(t, new(MastodonBridgeTestSuite))
}
