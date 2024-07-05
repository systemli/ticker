package bridge

import (
	"context"
	"errors"

	"github.com/mattn/go-mastodon"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type MastodonBridge struct {
	config  config.Config
	storage storage.Storage
}

func (mb *MastodonBridge) UpdateTicker(ticker storage.Ticker) error {
	return nil
}

func (mb *MastodonBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	if !ticker.Mastodon.Active {
		return nil
	}

	ctx := context.Background()
	client := client(ticker)

	var mediaIDs []mastodon.ID
	if len(message.Attachments) > 0 {
		for _, attachment := range message.Attachments {
			upload, err := mb.storage.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			media, err := client.UploadMedia(ctx, upload.FullPath(mb.config.Upload.Path))
			if err != nil {
				log.WithError(err).Error("unable to upload the attachment")
				continue
			}
			mediaIDs = append(mediaIDs, media.ID)
		}
	}

	toot := mastodon.Toot{
		Status:   message.Text,
		MediaIDs: mediaIDs,
	}

	status, err := client.PostStatus(ctx, &toot)
	if err != nil {
		return err
	}

	message.Mastodon = storage.MastodonMeta{
		ID:  string(status.ID),
		URI: status.URI,
		URL: status.URL,
	}

	return nil
}

func (mb *MastodonBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	if message.Mastodon.ID == "" {
		return nil
	}

	if !ticker.Mastodon.Connected() {
		return errors.New("unable to delete the status")
	}

	ctx := context.Background()
	client := client(ticker)

	return client.DeleteStatus(ctx, mastodon.ID(message.Mastodon.ID))
}

func client(ticker storage.Ticker) *mastodon.Client {
	return mastodon.NewClient(&mastodon.Config{
		Server:       ticker.Mastodon.Server,
		ClientID:     ticker.Mastodon.Token,
		ClientSecret: ticker.Mastodon.Secret,
		AccessToken:  ticker.Mastodon.AccessToken,
	})
}
