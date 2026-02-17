package bridge

import (
	"errors"

	"github.com/h2non/gock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func (s *BridgeTestSuite) TestSignalGroupUpdate() {
	s.Run("when signalGroup is inactive", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.DefaultSignalGroupSettings())
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		err := bridge.Update(tickerWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when signalGroup is active but signal-cli api fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := bridge.Update(tickerWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("happy path", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		// updateGroup
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "sample-group-id",
					"timestamp": 1,
				},
				"id": 1,
			})
		// listGroups
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "sample-group-id",
						"name":            "Example",
						"description":     "Example",
						"groupInviteLink": "https://signal.group/#example",
					},
				},
				"id": 1,
			})

		err := bridge.Update(tickerWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})
}

func (s *BridgeTestSuite) TestSignalGroupSend() {
	s.Run("when signalGroup is inactive", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.DefaultSignalGroupSettings())
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		err := bridge.Send(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when signalGroup is active but signal-cli api fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when response timestamp == 0", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]int{
					"timestamp": 0,
				},
				"id": 1,
			})

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("send message with attachment failed to find upload", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, errors.New("failed to find upload")).Once()
		mockStorage.On("FindUploadByUUID", "456").Return(storage.Upload{}, errors.New("failed to find upload")).Once()
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]int{
					"timestamp": 1,
				},
				"id": 1,
			})

		err := bridge.Send(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("send message with attachment failed to read file", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{UUID: "123", ContentType: "image/gif"}, nil).Once()
		mockStorage.On("FindUploadByUUID", "456").Return(storage.Upload{UUID: "456", ContentType: "image/jpeg"}, nil).Once()
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]int{
					"timestamp": 1,
				},
				"id": 1,
			})

		err := bridge.Send(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("send message without attachments", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]int{
					"timestamp": 1,
				},
				"id": 1,
			})

		err := bridge.Send(tickerWithBridges, &storage.Message{})
		s.NoError(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})
}

func (s *BridgeTestSuite) TestSignalDelete() {
	s.Run("when signal not connected", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.DefaultSignalGroupSettings())
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when message has no signal meta", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when signal is inactive", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.DefaultSignalGroupSettings())
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("when delete fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})

	s.Run("happy path", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("GetSignalGroupSettings").Return(storage.SignalGroupSettings{
			ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
			Account: "0123456789",
		})
		bridge := s.signalGroupBridge(config.Config{}, mockStorage)

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]int{
					"timestamp": 1,
				},
				"id": 1,
			})

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
		mockStorage.AssertExpectations(s.T())
	})
}

func (s *BridgeTestSuite) signalGroupBridge(config config.Config, storage storage.Storage) *SignalGroupBridge {
	return &SignalGroupBridge{
		config:  config,
		storage: storage,
	}
}
