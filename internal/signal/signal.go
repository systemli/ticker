package signal

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
	"github.com/ybbus/jsonrpc/v3"
)

var log = logrus.WithField("package", "signal")

type createGroupParams struct {
	Account                   string `json:"account"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	Avatar                    string `json:"avatar"`
	Link                      string `json:"link"`
	SetPermissionAddMember    string `json:"setPermissionAddMember"`
	SetPermissionEditDetails  string `json:"setPermissionEditDetails"`
	SetPermissionSendMessages string `json:"setPermissionSendMessages"`
	Expiration                int    `json:"expiration"`
}

type CreateGroupResponse struct {
	GroupID   string `json:"groupId"`
	Timestamp int    `json:"timestamp"`
}

type updateGroupParams struct {
	Account                   string `json:"account"`
	GroupID                   string `json:"group-id"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	Avatar                    string `json:"avatar"`
	Link                      string `json:"link"`
	SetPermissionAddMember    string `json:"setPermissionAddMember"`
	SetPermissionEditDetails  string `json:"setPermissionEditDetails"`
	SetPermissionSendMessages string `json:"setPermissionSendMessages"`
	Expiration                int    `json:"expiration"`
}

type UpdateGroupResponse struct {
	Timestamp int `json:"timestamp"`
}

type QuitGroupParams struct {
	Account string `json:"account"`
	GroupID string `json:"group-id"`
	Delete  bool   `json:"delete"`
}

type ListGroupsParams struct {
	Account string `json:"account"`
}

type ListGroupsResponseGroup struct {
	GroupID         string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	GroupInviteLink string `json:"groupInviteLink"`
}

type SendParams struct {
	Account    string   `json:"account"`
	GroupID    string   `json:"group-id"`
	Message    string   `json:"message"`
	Attachment []string `json:"attachment"`
}

type SendResponse struct {
	Timestamp *int `json:"timestamp"`
}

type DeleteParams struct {
	Account         string `json:"account"`
	GroupID         string `json:"group-id"`
	TargetTimestamp *int   `json:"target-timestamp"`
}

func CreateOrUpdateGroup(ts *storage.TickerSignalGroup, config config.Config) error {
	ctx := context.Background()
	client := rpcClient(config)

	var err error
	if ts.GroupID == "" {
		// Create new group
		var response *CreateGroupResponse
		params := createGroupParams{
			Account:                   config.SignalGroup.Account,
			Name:                      ts.GroupName,
			Description:               ts.GroupDescription,
			Avatar:                    "/var/lib/signal-cli/data/ticker.png",
			Link:                      "enabled",
			SetPermissionAddMember:    "every-member",
			SetPermissionEditDetails:  "only-admins",
			SetPermissionSendMessages: "only-admins",
			Expiration:                86400,
		}
		err = client.CallFor(ctx, &response, "updateGroup", &params)
		if err != nil {
			return err
		}
		if response.GroupID == "" {
			return errors.New("SignalGroup Bridge: No group ID in create group response")
		}
		log.WithField("groupId", response.GroupID).Debug("Created group")
		ts.GroupID = response.GroupID
	} else {
		// Update existing group
		params := updateGroupParams{
			Account:                   config.SignalGroup.Account,
			GroupID:                   ts.GroupID,
			Name:                      ts.GroupName,
			Description:               ts.GroupDescription,
			Avatar:                    "/var/lib/signal-cli/data/ticker.png",
			Link:                      "enabled",
			SetPermissionAddMember:    "every-member",
			SetPermissionEditDetails:  "only-admins",
			SetPermissionSendMessages: "only-admins",
			Expiration:                86400,
		}
		var response *UpdateGroupResponse
		err = client.CallFor(ctx, &response, "updateGroup", &params)
		if err != nil {
			return err
		}
	}

	g, err := getGroup(config, ts.GroupID)
	if err != nil {
		return err
	}
	if g == nil {
		return errors.New("SignalGroup Bridge: Group not found")
	}
	if g.GroupInviteLink == "" {
		return errors.New("SignalGroup Bridge: No invite link in group response")
	}

	ts.GroupInviteLink = g.GroupInviteLink

	return nil
}

func QuitGroup(config config.Config, groupID string) error {
	ctx := context.Background()
	client := rpcClient(config)

	params := QuitGroupParams{
		Account: config.SignalGroup.Account,
		GroupID: groupID,
		Delete:  true,
	}

	// TODO: cannot leave group if I'm the last admin
	// Maybe promote first other member to admin?
	var response interface{}
	err := client.CallFor(ctx, &response, "leaveGroup", &params)
	if err != nil {
		return err
	}

	return nil
}

func listGroups(config config.Config) ([]*ListGroupsResponseGroup, error) {
	ctx := context.Background()
	client := rpcClient(config)

	params := ListGroupsParams{
		Account: config.SignalGroup.Account,
	}

	var response []*ListGroupsResponseGroup
	err := client.CallFor(ctx, &response, "listGroups", &params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getGroup(config config.Config, groupID string) (*ListGroupsResponseGroup, error) {
	gl, err := listGroups(config)
	if err != nil {
		return nil, err
	}

	for _, g := range gl {
		if g.GroupID == groupID {
			return g, nil
		}
	}

	return nil, nil
}

func SendGroupMessage(config config.Config, ss storage.Storage, groupID string, message *storage.Message) error {
	ctx := context.Background()
	client := rpcClient(config)

	var attachments []string
	if len(message.Attachments) > 0 {
		for _, attachment := range message.Attachments {
			upload, err := ss.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			fileContent, err := os.ReadFile(upload.FullPath(config.Upload.Path))
			if err != nil {
				log.WithError(err).Error("failed to read file")
				continue
			}
			fileBase64 := base64.StdEncoding.EncodeToString(fileContent)
			aString := fmt.Sprintf("data:%s;filename=%s;base64,%s", upload.ContentType, upload.FileName, fileBase64)
			attachments = append(attachments, aString)
		}
	}

	params := SendParams{
		Account:    config.SignalGroup.Account,
		GroupID:    groupID,
		Message:    message.Text,
		Attachment: attachments,
	}

	var response *SendResponse
	err := client.CallFor(ctx, &response, "send", &params)
	if err != nil {
		return err
	}
	if response.Timestamp == nil {
		return errors.New("SignalGroup Bridge: No timestamp in send response")
	}

	message.SignalGroup = storage.SignalGroupMeta{
		Timestamp: response.Timestamp,
	}

	return nil
}

func DeleteMessage(config config.Config, groupID string, message *storage.Message) error {
	ctx := context.Background()
	client := rpcClient(config)

	params := DeleteParams{
		Account:         config.SignalGroup.Account,
		GroupID:         groupID,
		TargetTimestamp: message.SignalGroup.Timestamp,
	}

	var response *SendResponse
	err := client.CallFor(ctx, &response, "remoteDelete", &params)
	if err != nil {
		return err
	}

	return nil
}

func rpcClient(config config.Config) jsonrpc.RPCClient {
	if config.SignalGroup.ApiUser != "" && config.SignalGroup.ApiPass != "" {
		return jsonrpc.NewClientWithOpts(config.SignalGroup.ApiUrl, &jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.SignalGroup.ApiUser+":"+config.SignalGroup.ApiPass)),
			},
		})
	} else {
		return jsonrpc.NewClient(config.SignalGroup.ApiUrl)

	}
}
