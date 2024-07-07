package bridge

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/signal"
	"github.com/systemli/ticker/internal/storage"
)

type SignalGroupBridge struct {
	config  config.Config
	storage storage.Storage
}

type SignalGroupResponse struct {
	Timestamp int `json:"timestamp"`
}

func (sb *SignalGroupBridge) Update(ticker storage.Ticker) error {
	if !sb.config.SignalGroup.Enabled() || !ticker.SignalGroup.Connected() {
		return nil
	}

	groupClient := signal.NewGroupClient(sb.config)
	err := groupClient.CreateOrUpdateGroup(&ticker)
	if err != nil {
		return err
	}

	return nil
}

func (sb *SignalGroupBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	if !sb.config.SignalGroup.Enabled() || !ticker.SignalGroup.Connected() || !ticker.SignalGroup.Active {
		return nil
	}

	ctx := context.Background()
	client := signal.Client(sb.config)

	var attachments []string
	if len(message.Attachments) > 0 {
		for _, attachment := range message.Attachments {
			upload, err := sb.storage.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			fileContent, err := os.ReadFile(upload.FullPath(sb.config.Upload.Path))
			if err != nil {
				log.WithError(err).Error("failed to read file")
				continue
			}
			fileBase64 := base64.StdEncoding.EncodeToString(fileContent)
			aString := fmt.Sprintf("data:%s;filename=%s;base64,%s", upload.ContentType, upload.FileName(), fileBase64)
			attachments = append(attachments, aString)
		}
	}

	params := struct {
		Account    string   `json:"account"`
		GroupID    string   `json:"group-id"`
		Message    string   `json:"message"`
		Attachment []string `json:"attachment"`
	}{
		Account:    sb.config.SignalGroup.Account,
		GroupID:    ticker.SignalGroup.GroupID,
		Message:    message.Text,
		Attachment: attachments,
	}

	var response SignalGroupResponse
	err := client.CallFor(ctx, &response, "send", &params)
	if err != nil {
		return err
	}
	if response.Timestamp == 0 {
		return errors.New("SignalGroup Bridge: No timestamp in send response")
	}

	message.SignalGroup = storage.SignalGroupMeta{
		Timestamp: response.Timestamp,
	}

	return nil
}

func (sb *SignalGroupBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	if !sb.config.SignalGroup.Enabled() || !ticker.SignalGroup.Connected() || !ticker.SignalGroup.Active || message.SignalGroup.Timestamp == 0 {
		return nil
	}

	client := signal.Client(sb.config)
	params := struct {
		Account         string `json:"account"`
		GroupID         string `json:"group-id"`
		TargetTimestamp int    `json:"target-timestamp"`
	}{
		Account:         sb.config.SignalGroup.Account,
		GroupID:         ticker.SignalGroup.GroupID,
		TargetTimestamp: message.SignalGroup.Timestamp,
	}

	var response SignalGroupResponse
	err := client.CallFor(context.Background(), &response, "remoteDelete", &params)
	if err != nil {
		return err
	}

	return nil
}
