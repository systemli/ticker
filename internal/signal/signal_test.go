package signal

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type GroupClientTestSuite struct {
	suite.Suite
	settings storage.SignalGroupSettings
}

func (s *GroupClientTestSuite) SetupTest() {
	gock.DisableNetworking()
	s.settings = storage.SignalGroupSettings{
		ApiUrl:  "https://signal-cli.example.org/api/v1/rpc",
		Account: "+491234567890",
		Avatar:  "/path/to/avatar.png",
	}
}

func (s *GroupClientTestSuite) TearDownTest() {
	gock.Off()
}

func (s *GroupClientTestSuite) newClient() *GroupClient {
	return NewGroupClientFromSettings(s.settings)
}

func (s *GroupClientTestSuite) TestNewGroupClientFromSettings() {
	gc := s.newClient()
	s.Equal(s.settings, gc.settings)
	s.NotNil(gc.client)
}

func (s *GroupClientTestSuite) TestCreateOrUpdateGroup() {
	s.Run("creates a new group", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Test Ticker",
			Description: "Test Description",
			SignalGroup: storage.TickerSignalGroup{},
		}

		// updateGroup call
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "new-group-id",
					"timestamp": 1,
				},
				"id": 1,
			})
		// listGroups call
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "new-group-id",
						"name":            "Test Ticker",
						"description":     "Test Description",
						"groupInviteLink": "https://signal.group/#example-link",
					},
				},
				"id": 2,
			})

		err := gc.CreateOrUpdateGroup(ticker)
		s.NoError(err)
		s.Equal("new-group-id", ticker.SignalGroup.GroupID)
		s.Equal("https://signal.group/#example-link", ticker.SignalGroup.GroupInviteLink)
		s.True(gock.IsDone())
	})

	s.Run("updates an existing group", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Updated Ticker",
			Description: "Updated Description",
			SignalGroup: storage.TickerSignalGroup{
				GroupID: "existing-group-id",
			},
		}

		// updateGroup call
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "existing-group-id",
					"timestamp": 2,
				},
				"id": 1,
			})
		// listGroups call
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "existing-group-id",
						"name":            "Updated Ticker",
						"description":     "Updated Description",
						"groupInviteLink": "https://signal.group/#updated-link",
					},
				},
				"id": 2,
			})

		err := gc.CreateOrUpdateGroup(ticker)
		s.NoError(err)
		s.Equal("existing-group-id", ticker.SignalGroup.GroupID)
		s.Equal("https://signal.group/#updated-link", ticker.SignalGroup.GroupInviteLink)
		s.True(gock.IsDone())
	})

	s.Run("returns error when updateGroup RPC fails", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Test Ticker",
			Description: "Test Description",
			SignalGroup: storage.TickerSignalGroup{},
		}

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.CreateOrUpdateGroup(ticker)
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error when response has no groupId and ticker has no groupId", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Test Ticker",
			Description: "Test Description",
			SignalGroup: storage.TickerSignalGroup{},
		}

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "",
					"timestamp": 1,
				},
				"id": 1,
			})

		err := gc.CreateOrUpdateGroup(ticker)
		s.Error(err)
		s.Equal("unable to create or update group", err.Error())
		s.True(gock.IsDone())
	})

	s.Run("returns error when listGroups fails", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Test Ticker",
			Description: "Test Description",
			SignalGroup: storage.TickerSignalGroup{
				GroupID: "existing-group-id",
			},
		}

		// updateGroup succeeds
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "existing-group-id",
					"timestamp": 1,
				},
				"id": 1,
			})
		// listGroups fails
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.CreateOrUpdateGroup(ticker)
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error when group has no invite link", func() {
		gc := s.newClient()
		ticker := &storage.Ticker{
			Title:       "Test Ticker",
			Description: "Test Description",
			SignalGroup: storage.TickerSignalGroup{
				GroupID: "existing-group-id",
			},
		}

		// updateGroup succeeds
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"groupId":   "existing-group-id",
					"timestamp": 1,
				},
				"id": 1,
			})
		// listGroups returns group without invite link
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "existing-group-id",
						"name":            "Test Ticker",
						"description":     "Test Description",
						"groupInviteLink": "",
					},
				},
				"id": 2,
			})

		err := gc.CreateOrUpdateGroup(ticker)
		s.Error(err)
		s.Equal("unable to get group invite link", err.Error())
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestQuitGroup() {
	s.Run("happy path", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  map[string]int{"timestamp": 1},
				"id":      1,
			})

		err := gc.QuitGroup("group-id-123")
		s.NoError(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error on RPC failure", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.QuitGroup("group-id-123")
		s.Error(err)
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestListGroups() {
	s.Run("happy path", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "group-1",
						"name":            "Group One",
						"description":     "First group",
						"groupInviteLink": "https://signal.group/#link1",
						"members": []map[string]string{
							{"number": "+491111111111", "uuid": "uuid-1"},
						},
					},
					{
						"id":              "group-2",
						"name":            "Group Two",
						"description":     "Second group",
						"groupInviteLink": "https://signal.group/#link2",
						"members":         []map[string]string{},
					},
				},
				"id": 1,
			})

		groups, err := gc.ListGroups()
		s.NoError(err)
		s.Len(groups, 2)
		s.Equal("group-1", groups[0].GroupID)
		s.Equal("Group One", groups[0].Name)
		s.Equal("First group", groups[0].Description)
		s.Equal("https://signal.group/#link1", groups[0].GroupInviteLink)
		s.Len(groups[0].Members, 1)
		s.Equal("+491111111111", groups[0].Members[0].Number)
		s.Equal("uuid-1", groups[0].Members[0].Uuid)
		s.Equal("group-2", groups[1].GroupID)
		s.Equal("Group Two", groups[1].Name)
		s.True(gock.IsDone())
	})

	s.Run("returns empty list", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  []map[string]interface{}{},
				"id":      1,
			})

		groups, err := gc.ListGroups()
		s.NoError(err)
		s.Empty(groups)
		s.True(gock.IsDone())
	})

	s.Run("returns error on RPC failure", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		groups, err := gc.ListGroups()
		s.Error(err)
		s.Nil(groups)
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestAddAdminMember() {
	s.Run("happy path", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  map[string]int{"timestamp": 1},
				"id":      1,
			})

		err := gc.AddAdminMember("group-id-123", "+499876543210")
		s.NoError(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error on RPC failure", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.AddAdminMember("group-id-123", "+499876543210")
		s.Error(err)
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestRemoveAllMembers() {
	s.Run("removes all members except account", func() {
		gc := s.newClient()

		// listGroups
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":   "group-id-123",
						"name": "Test Group",
						"members": []map[string]string{
							{"number": "+491234567890", "uuid": "uuid-self"},
							{"number": "+491111111111", "uuid": "uuid-1"},
							{"number": "+492222222222", "uuid": "uuid-2"},
						},
						"groupInviteLink": "https://signal.group/#link",
					},
				},
				"id": 1,
			})
		// updateGroup (removeMembers)
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  map[string]int{"timestamp": 1},
				"id":      2,
			})

		err := gc.RemoveAllMembers("group-id-123")
		s.NoError(err)
		s.True(gock.IsDone())
	})

	s.Run("does nothing when only the account is a member", func() {
		gc := s.newClient()

		// listGroups
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":   "group-id-123",
						"name": "Test Group",
						"members": []map[string]string{
							{"number": "+491234567890", "uuid": "uuid-self"},
						},
						"groupInviteLink": "https://signal.group/#link",
					},
				},
				"id": 1,
			})

		err := gc.RemoveAllMembers("group-id-123")
		s.NoError(err)
		s.True(gock.IsDone())
	})

	s.Run("does nothing when group has no members", func() {
		gc := s.newClient()

		// listGroups
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "group-id-123",
						"name":            "Test Group",
						"members":         []map[string]string{},
						"groupInviteLink": "https://signal.group/#link",
					},
				},
				"id": 1,
			})

		err := gc.RemoveAllMembers("group-id-123")
		s.NoError(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error when listGroups fails", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.RemoveAllMembers("group-id-123")
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("returns error when removeMembers fails", func() {
		gc := s.newClient()

		// listGroups succeeds
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":   "group-id-123",
						"name": "Test Group",
						"members": []map[string]string{
							{"number": "+491234567890", "uuid": "uuid-self"},
							{"number": "+491111111111", "uuid": "uuid-1"},
						},
						"groupInviteLink": "https://signal.group/#link",
					},
				},
				"id": 1,
			})
		// updateGroup (removeMembers) fails
		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		err := gc.RemoveAllMembers("group-id-123")
		s.Error(err)
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestGetGroup() {
	s.Run("finds existing group", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "group-a",
						"name":            "Group A",
						"groupInviteLink": "https://signal.group/#a",
					},
					{
						"id":              "group-b",
						"name":            "Group B",
						"groupInviteLink": "https://signal.group/#b",
					},
				},
				"id": 1,
			})

		g, err := gc.getGroup("group-b")
		s.NoError(err)
		s.Equal("group-b", g.GroupID)
		s.Equal("Group B", g.Name)
		s.Equal("https://signal.group/#b", g.GroupInviteLink)
		s.True(gock.IsDone())
	})

	s.Run("returns empty struct when group not found", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"id":              "group-a",
						"name":            "Group A",
						"groupInviteLink": "https://signal.group/#a",
					},
				},
				"id": 1,
			})

		g, err := gc.getGroup("nonexistent")
		s.NoError(err)
		s.Equal(ListGroupsResponseGroup{}, g)
		s.True(gock.IsDone())
	})

	s.Run("returns error when listGroups fails", func() {
		gc := s.newClient()

		gock.New("https://signal-cli.example.org").
			Post("/api/v1/rpc").
			Reply(500)

		g, err := gc.getGroup("group-a")
		s.Error(err)
		s.Equal(ListGroupsResponseGroup{}, g)
		s.True(gock.IsDone())
	})
}

func (s *GroupClientTestSuite) TestClientFromSettings() {
	client := ClientFromSettings(s.settings)
	s.NotNil(client)
}

func TestGroupClientTestSuite(t *testing.T) {
	suite.Run(t, new(GroupClientTestSuite))
}
