package bridge

import (
	"errors"

	"github.com/h2non/gock"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func (s *BridgeTestSuite) TestBlueskyUpdate() {
	s.Run("does nothing", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Update(tickerWithBridges)
		s.NoError(err)
	})
}

func (s *BridgeTestSuite) TestBlueskySend() {
	s.Run("when bluesky is inactive", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Send(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when bluesky is active but login fails", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(401)

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("when bluesky is active and login succeeds", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.blueskyBridge(config.Config{}, mockStorage)

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(200).
			JSON(map[string]string{
				"uri": "sample-uri",
				"cid": "sample-cid",
			})

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		s.Equal("sample-uri", messageWithoutBridges.Bluesky.Uri)
		s.Equal("sample-cid", messageWithoutBridges.Bluesky.Cid)
		s.Equal("handle", messageWithoutBridges.Bluesky.Handle)

		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when bluesky is active and upload is not found", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, errors.New("not found")).Once()
		bridge := s.blueskyBridge(config.Config{}, mockStorage)

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(200).
			JSON(map[string]string{
				"uri": "sample-uri",
				"cid": "sample-cid",
			})

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
		s.Equal("sample-uri", messageWithoutBridges.Bluesky.Uri)
		s.Equal("sample-cid", messageWithoutBridges.Bluesky.Cid)
		s.Equal("handle", messageWithoutBridges.Bluesky.Handle)

		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when bluesky is active but bluesky responds with error", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.blueskyBridge(config.Config{}, mockStorage)

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(500)

		err := bridge.Send(tickerWithBridges, &messageWithoutBridges)
		s.Error(err)
		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestBlueskyDelete() {
	s.Run("when bluesky not connected", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithoutBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when message has no bluesky meta", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when bluesky is inactive", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		err := bridge.Delete(tickerWithBridges, &messageWithoutBridges)
		s.NoError(err)
	})

	s.Run("when bluesky is active but login fails", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(401)

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("when delete fails", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		// Thread gate deletion (may or may not exist)
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(200).
			JSON(map[string]string{})

		// Post deletion fails
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(500)

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("happy path", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		// Thread gate deletion (may or may not exist)
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(200).
			JSON(map[string]string{})

		// Post deletion
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(200).
			JSON(map[string]string{})

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
	})
}

func (s *BridgeTestSuite) blueskyBridge(config config.Config, storage storage.Storage) *BlueskyBridge {
	return &BlueskyBridge{
		config:  config,
		storage: storage,
	}
}

func (s *BridgeTestSuite) TestBuildAllowRules() {
	s.Run("nobody", func() {
		rules := buildAllowRules("nobody")
		s.NotNil(rules)
		s.Len(rules, 0)
	})

	s.Run("followers", func() {
		rules := buildAllowRules("followers")
		s.Len(rules, 1)
		s.NotNil(rules[0].FeedThreadgate_FollowerRule)
	})

	s.Run("following", func() {
		rules := buildAllowRules("following")
		s.Len(rules, 1)
		s.NotNil(rules[0].FeedThreadgate_FollowingRule)
	})

	s.Run("mentioned", func() {
		rules := buildAllowRules("mentioned")
		s.Len(rules, 1)
		s.NotNil(rules[0].FeedThreadgate_MentionRule)
	})

	s.Run("followers,mentioned", func() {
		rules := buildAllowRules("followers,mentioned")
		s.Len(rules, 2)
		s.NotNil(rules[0].FeedThreadgate_FollowerRule)
		s.NotNil(rules[1].FeedThreadgate_MentionRule)
	})

	s.Run("followers,following,mentioned", func() {
		rules := buildAllowRules("followers,following,mentioned")
		s.Len(rules, 3)
		s.NotNil(rules[0].FeedThreadgate_FollowerRule)
		s.NotNil(rules[1].FeedThreadgate_FollowingRule)
		s.NotNil(rules[2].FeedThreadgate_MentionRule)
	})

	s.Run("empty string", func() {
		rules := buildAllowRules("")
		s.Nil(rules)
	})
}

func (s *BridgeTestSuite) TestBlueskySendWithThreadGate() {
	s.Run("when reply restriction is set to followers", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.blueskyBridge(config.Config{}, mockStorage)

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(200).
			JSON(map[string]string{
				"uri": "at://did:plc:sample-did/app.bsky.feed.post/sample-rkey",
				"cid": "sample-cid",
			})

		// Thread gate creation
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(200).
			JSON(map[string]string{
				"uri": "at://did:plc:sample-did/app.bsky.feed.threadgate/sample-rkey",
				"cid": "sample-cid-tg",
			})

		tickerWithReplyRestriction := storage.Ticker{
			Bluesky: storage.TickerBluesky{
				Active:           true,
				Handle:           "handle",
				AppKey:           "app_key",
				ReplyRestriction: "followers",
			},
		}

		msg := storage.Message{
			Text: "Hello World https://example.com",
			Attachments: []storage.Attachment{
				{UUID: "123"},
			},
		}

		err := bridge.Send(tickerWithReplyRestriction, &msg)
		s.NoError(err)
		s.Equal("at://did:plc:sample-did/app.bsky.feed.post/sample-rkey", msg.Bluesky.Uri)
		s.Equal("sample-cid", msg.Bluesky.Cid)
		s.Equal("handle", msg.Bluesky.Handle)

		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})

	s.Run("when thread gate creation fails", func() {
		mockStorage := &storage.MockStorage{}
		mockStorage.On("FindUploadByUUID", "123").Return(storage.Upload{}, nil).Once()
		bridge := s.blueskyBridge(config.Config{}, mockStorage)

		gock.DisableNetworking()
		defer gock.Off()

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(200).
			JSON(map[string]string{
				"uri": "at://did:plc:sample-did/app.bsky.feed.post/sample-rkey",
				"cid": "sample-cid",
			})

		// Thread gate creation fails
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.createRecord").
			Reply(500)

		tickerWithReplyRestriction := storage.Ticker{
			Bluesky: storage.TickerBluesky{
				Active:           true,
				Handle:           "handle",
				AppKey:           "app_key",
				ReplyRestriction: "nobody",
			},
		}

		msg := storage.Message{
			Text: "Hello World https://example.com",
			Attachments: []storage.Attachment{
				{UUID: "123"},
			},
		}

		// Should still succeed - thread gate failure is non-fatal
		err := bridge.Send(tickerWithReplyRestriction, &msg)
		s.NoError(err)
		s.Equal("at://did:plc:sample-did/app.bsky.feed.post/sample-rkey", msg.Bluesky.Uri)

		s.True(gock.IsDone())
		s.True(mockStorage.AssertExpectations(s.T()))
	})
}

func (s *BridgeTestSuite) TestBlueskyDeleteWithThreadGate() {
	s.Run("happy path with thread gate deletion", func() {
		bridge := s.blueskyBridge(config.Config{}, &storage.MockStorage{})

		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.server.createSession").
			Reply(200).
			JSON(map[string]string{
				"Did":        "sample-did",
				"AccessJwt":  "sample-access-jwt",
				"RefreshJwt": "sample-refresh-jwt",
			})

		// Thread gate deletion
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(200).
			JSON(map[string]string{})

		// Post deletion
		gock.New("https://bsky.social").
			Post("/xrpc/com.atproto.repo.deleteRecord").
			Reply(200).
			JSON(map[string]string{})

		err := bridge.Delete(tickerWithBridges, &messageWithBridges)
		s.NoError(err)
		s.True(gock.IsDone())
	})
}
