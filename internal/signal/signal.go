package signal

import (
	"context"
	"errors"

	"github.com/systemli/ticker/internal/storage"
	"github.com/ybbus/jsonrpc/v3"
)

type GroupMember struct {
	Number string `json:"number"`
	Uuid   string `json:"uuid"`
}

type ListGroupsResponseGroup struct {
	GroupID         string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Members         []GroupMember `json:"members"`
	GroupInviteLink string        `json:"groupInviteLink"`
}

type GroupClient struct {
	settings storage.SignalGroupSettings
	client   jsonrpc.RPCClient
}

func NewGroupClientFromSettings(settings storage.SignalGroupSettings) *GroupClient {
	client := ClientFromSettings(settings)
	return &GroupClient{settings, client}
}

func (gc *GroupClient) CreateOrUpdateGroup(t *storage.Ticker) error {
	params := map[string]interface{}{
		"account":                   gc.settings.Account,
		"name":                      t.Title,
		"description":               t.Description,
		"avatar":                    gc.settings.Avatar,
		"link":                      "enabled",
		"setPermissionAddMember":    "every-member",
		"setPermissionEditDetails":  "only-admins",
		"setPermissionSendMessages": "only-admins",
		"expiration":                86400,
	}
	if t.SignalGroup.GroupID != "" {
		params["group-id"] = t.SignalGroup.GroupID
	}

	var response struct {
		GroupID   string `json:"groupId"`
		Timestamp int    `json:"timestamp"`
	}
	err := gc.client.CallFor(context.Background(), &response, "updateGroup", &params)
	if err != nil {
		return err
	}
	if t.SignalGroup.GroupID == "" {
		// Set GroupID for newly created group
		t.SignalGroup.GroupID = response.GroupID
	}

	if t.SignalGroup.GroupID == "" {
		return errors.New("unable to create or update group")
	}

	g, err := gc.getGroup(t.SignalGroup.GroupID)
	if err != nil {
		return err
	}
	if g.GroupInviteLink == "" {
		return errors.New("unable to get group invite link")
	}

	t.SignalGroup.GroupInviteLink = g.GroupInviteLink

	return nil
}

func (gc *GroupClient) QuitGroup(groupID string) error {
	params := struct {
		Account string `json:"account"`
		GroupID string `json:"group-id"`
		Delete  bool   `json:"delete"`
	}{
		Account: gc.settings.Account,
		GroupID: groupID,
		Delete:  true,
	}

	var response interface{}
	err := gc.client.CallFor(context.Background(), &response, "quitGroup", &params)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GroupClient) ListGroups() ([]ListGroupsResponseGroup, error) {
	ctx := context.Background()

	params := struct {
		Account  string `json:"account"`
		Detailed bool   `json:"detailed"`
	}{
		Account:  gc.settings.Account,
		Detailed: true,
	}

	var response []ListGroupsResponseGroup
	err := gc.client.CallFor(ctx, &response, "listGroups", &params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (gc *GroupClient) getGroup(groupID string) (ListGroupsResponseGroup, error) {
	gl, err := gc.ListGroups()
	if err != nil {
		return ListGroupsResponseGroup{}, err
	}

	for _, g := range gl {
		if g.GroupID == groupID {
			return g, nil
		}
	}

	return ListGroupsResponseGroup{}, nil
}

func (gc *GroupClient) AddAdminMember(groupId string, number string) error {
	numbers := make([]string, 0, 1)
	numbers = append(numbers, number)

	params := struct {
		Account string   `json:"account"`
		GroupID string   `json:"group-id"`
		Member  []string `json:"member"`
		Admin   []string `json:"admin"`
	}{
		Account: gc.settings.Account,
		GroupID: groupId,
		Member:  numbers,
		Admin:   numbers,
	}

	var response interface{}
	err := gc.client.CallFor(context.Background(), &response, "updateGroup", &params)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GroupClient) RemoveAllMembers(groupId string) error {
	g, err := gc.getGroup(groupId)
	if err != nil {
		return err
	}

	numbers := make([]string, 0, len(g.Members))
	for _, m := range g.Members {
		// Exclude the account number
		if m.Number == gc.settings.Account {
			continue
		}
		numbers = append(numbers, m.Number)
	}

	if len(numbers) == 0 {
		return nil
	}

	return gc.removeMembers(groupId, numbers)
}

func (gc *GroupClient) removeMembers(groupId string, numbers []string) error {
	params := struct {
		Account      string   `json:"account"`
		GroupID      string   `json:"group-id"`
		RemoveMember []string `json:"remove-member"`
	}{
		Account:      gc.settings.Account,
		GroupID:      groupId,
		RemoveMember: numbers,
	}

	var response interface{}
	err := gc.client.CallFor(context.Background(), &response, "updateGroup", &params)
	if err != nil {
		return err
	}

	return nil
}

func ClientFromSettings(settings storage.SignalGroupSettings) jsonrpc.RPCClient {
	return jsonrpc.NewClient(settings.ApiUrl)
}
