package bridge

import (
	"io/ioutil"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func SendTweet(ticker *model.Ticker, message *model.Message) error {
	if !ticker.Twitter.Active {
		return nil
	}

	if err := TwitterConnectionEnabled(ticker); err != nil {
		return err
	}

	client := TwitterClient(ticker.Twitter.Token, ticker.Twitter.Secret)
	params := &twitter.StatusUpdateParams{}

	if len(message.Attachments) > 0 {
		var mediaIds []int64
		for _, attachment := range message.Attachments {
			upload := &model.Upload{}
			err := storage.DB.One("UUID", attachment.UUID, upload)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}
			bytes, err := ioutil.ReadFile(upload.FullPath())
			if err != nil {
				log.WithError(err).Error("failed to open upload")
				continue
			}

			mpr, _, err := client.Media.Upload(bytes, upload.ContentType)
			if err != nil {
				log.WithError(err).Error("failed to upload the media file to twitter")
				continue
			}

			mediaIds = append(mediaIds, mpr.MediaID)
		}
		params.MediaIds = mediaIds
	}

	tweet, _, err := client.Statuses.Update(message.Text, params)
	if err != nil {
		return err
	}

	message.Tweet = model.Tweet{ID: tweet.IDStr, UserName: tweet.User.ScreenName}

	return nil
}

func DeleteTweet(ticker *model.Ticker, message *model.Message) error {
	if err := TwitterConnectionEnabled(ticker); err != nil {
		return err
	}

	if message.Tweet.ID == "" {
		return nil
	}

	id, err := strconv.ParseInt(message.Tweet.ID, 10, 64)
	if err != nil {
		return err
	}

	client := TwitterClient(ticker.Twitter.Token, ticker.Twitter.Secret)
	_, _, err = client.Statuses.Destroy(id, nil)

	return err
}

func TwitterUser(ticker *model.Ticker) (*twitter.User, error) {
	u := &twitter.User{}

	if err := TwitterConnectionEnabled(ticker); err != nil {
		return u, err
	}

	client := TwitterClient(ticker.Twitter.Token, ticker.Twitter.Secret)
	avp := &twitter.AccountVerifyParams{
		IncludeEmail:    twitter.Bool(false),
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(avp)
	if err != nil {
		return user, err
	}

	return user, nil
}
