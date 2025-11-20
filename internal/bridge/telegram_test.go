package bridge

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/h2non/gock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func (s *BridgeTestSuite) TestTelegramUpdate() {
	s.Run("does nothing", func() {
		bridge := s.telegramBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Update(tickerWithBridges)
		s.NoError(err)
	})
}

func (s *BridgeTestSuite) TestTelegramSend() {
	s.Run("when telegram is inactive", func() {
		mockStorage := &storage.MockStorage{}
		// No expectation for GetTelegramSettings since ChannelName is empty
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		err := bridge.Send(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active but token is empty", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: ""})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active but bot api fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active but send message fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/sendMessage").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("when telegram is active without attachments", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/sendMessage").
			Reply(200).
			JSON(map[string]interface{}{
				"ok": true,
				"result": map[string]interface{}{
					"message_id": 123,
				},
			})

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.NoError(err)
		s.Equal(123, messageWithBridges.Telegram.Messages[0].MessageID)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active with attachments", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{UUID: "123", ContentType: "image/gif"}, nil).Once()
		mockStorage.On("FindUploadByUUID", "456").Return(storage.Upload{UUID: "456", ContentType: "image/jpeg"}, nil).Once()
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/sendMediaGroup").
			Reply(200).
			JSON(map[string]interface{}{
				"ok": true,
				"result": []interface{}{
					map[string]interface{}{
						"message_id": 123,
					},
				},
			})

		err := bridge.Send(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when telegram is active but send media group fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{UUID: "123", ContentType: "image/gif"}, nil).Once()
		mockStorage.On("FindUploadByUUID", "456").Return(storage.Upload{UUID: "456", ContentType: "image/jpeg"}, nil).Once()
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/sendMediaGroup").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when telegram is active but find upload fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, errors.New("failed to find upload")).Once()
		mockStorage.On("FindUploadByUUID", "456").Return(storage.Upload{}, errors.New("failed to find upload")).Once()
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		err := bridge.Send(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestTelegramDelete() {
	s.Run("when telegram is inactive", func() {
		mockStorage := &storage.MockStorage{}
		// No expectation for GetTelegramSettings since ChannelName is empty
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when token is empty", func() {
		mockStorage := &storage.MockStorage{}
		// No GetTelegramSettings expectation because ChannelName is empty, causing early return
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when message has no telegram meta", func() {
		mockStorage := &storage.MockStorage{}
		// No GetTelegramSettings expectation because message has no Telegram metadata, causing early return
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active but bot api fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(500)

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when telegram is active but delete message fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/deleteMessage").
			Reply(500)

		_ = bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.True(gock.IsDone())
	})

	s.Run("when telegram is active and deletion is successful", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetTelegramSettings").Return(storage.TelegramSettings{Token: "123"})
		bridge := s.telegramBridge(config.Config{}, mockStorage)

		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		gock.New("https://api.telegram.org").
			Post("/bot123/deleteMessage").
			Reply(200).
			JSON(map[string]interface{}{
				"ok": true,
			})

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})
}

func (s *BridgeTestSuite) TestBotUser() {
	s.Run("when bot api fails", func() {
		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(500)

		user, err := BotUser("123")
		s.Error(err)
		s.Equal(tgbotapi.User{}, user)
		s.True(gock.IsDone())
	})

	s.Run("when bot api is successful", func() {
		gock.New("https://api.telegram.org").
			Post("/bot123/getMe").
			Reply(200).
			JSON(map[string]interface{}{
				"ok":     true,
				"result": map[string]interface{}{"id": 123},
			})

		user, err := BotUser("123")
		s.NoError(err)
		s.Equal(tgbotapi.User{ID: 123}, user)
		s.True(gock.IsDone())
	})
}

func (s *BridgeTestSuite) telegramBridge(config config.Config, storage storage.Storage) *TelegramBridge {
	return &TelegramBridge{
		config:  config,
		storage: storage,
	}
}
