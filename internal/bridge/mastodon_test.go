package bridge

import (
	"errors"

	"github.com/h2non/gock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func (s *BridgeTestSuite) TestMastodonUpdate() {
	s.Run("does nothing", func() {
		bridge := s.mastodonBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Update(tickerWithBridges)
		s.NoError(err)
	})
}

func (s *BridgeTestSuite) TestMastodonSend() {
	s.Run("when mastodon is inactive", func() {
		bridge := s.mastodonBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Send(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when mastodon is active but upload cant not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.mastodonBridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(200).
			JSON(map[string]string{
				"id":  "123",
				"uri": "https://systemli.social/@systemli/123",
				"url": "https://systemli.social/@systemli/123",
			})

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		s.Equal("123", messageWithoutBridges.Mastodon.ID)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutBridges.Mastodon.URI)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutBridges.Mastodon.URL)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when mastodon is active and upload is not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, errors.New("upload not found")).Once()
		bridge := s.mastodonBridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(200).
			JSON(map[string]string{
				"id":  "123",
				"uri": "https://systemli.social/@systemli/123",
				"url": "https://systemli.social/@systemli/123",
			})

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		s.Equal("123", messageWithoutBridges.Mastodon.ID)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutBridges.Mastodon.URI)
		s.Equal("https://systemli.social/@systemli/123", messageWithoutBridges.Mastodon.URL)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when mastodon is active but post status fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.mastodonBridge(config.Config{}, mockStorage)

		gock.New("https://systemli.social").
			Post("/api/v1/statuses").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.Error(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestMastodonDelete() {
	s.Run("when message has no mastodon meta", func() {
		bridge := s.mastodonBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when mastodon is inactive", func() {
		bridge := s.mastodonBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithoutBridges, &messageWithBridges)
		s.Error(err)
		s.Equal("unable to delete the status", err.Error())
	})

	s.Run("when mastodon is active", func() {
		bridge := s.mastodonBridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://systemli.social").
			Delete("/api/v1/statuses/123").
			Reply(200)

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
	})
}

func (s *BridgeTestSuite) mastodonBridge(config config.Config, storage storage.Storage) *MastodonBridge {
	return &MastodonBridge{
		config:  config,
		storage: storage,
	}
}
