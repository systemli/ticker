package bridge

import (
	"errors"
	"os"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

type TwitterBridge struct {
	config  config.Config
	storage storage.TickerStorage
}

func (tb *TwitterBridge) Send(ticker storage.Ticker, message storage.Message) error {
	if !ticker.Twitter.Active {
		return nil
	}

	if err := twitterConnectionEnabled(ticker, tb.config); err != nil {
		return nil
	}

	client := tb.TwitterClient(ticker.Twitter.Token, ticker.Twitter.Secret)
	params := &twitter.StatusUpdateParams{}

	if len(message.Attachments) > 0 {
		var mediaIds []int64
		for _, attachment := range message.Attachments {
			upload, err := tb.storage.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}
			bytes, err := os.ReadFile(upload.FullPath(tb.config.UploadPath))
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

	message.Tweet = storage.Tweet{ID: tweet.IDStr, UserName: tweet.User.ScreenName}

	return nil
}

func (tb *TwitterBridge) Delete(ticker storage.Ticker, message storage.Message) error {
	if err := twitterConnectionEnabled(ticker, tb.config); err != nil {
		return nil
	}

	if message.Tweet.ID == "" {
		return nil
	}

	id, err := strconv.ParseInt(message.Tweet.ID, 10, 64)
	if err != nil {
		return err
	}

	client := tb.TwitterClient(ticker.Twitter.Token, ticker.Twitter.Secret)
	_, _, err = client.Statuses.Destroy(id, nil)

	return err
}

func (tb *TwitterBridge) TwitterClient(t, s string) *twitter.Client {
	config := oauth1.NewConfig(tb.config.TwitterConsumerKey, tb.config.TwitterConsumerSecret)
	token := oauth1.NewToken(t, s)

	return twitter.NewClient(config.Client(oauth1.NoContext, token))
}

func TwitterUser(ticker storage.Ticker, config config.Config) (*twitter.User, error) {
	if err := twitterConnectionEnabled(ticker, config); err != nil {
		return &twitter.User{}, err
	}

	client := twitterClient(config, ticker.Twitter.Token, ticker.Twitter.Secret)
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

func twitterClient(config config.Config, token, secret string) *twitter.Client {
	oauthConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	oauthToken := oauth1.NewToken(token, secret)

	return twitter.NewClient(oauthConfig.Client(oauth1.NoContext, oauthToken))
}

func twitterConnectionEnabled(ticker storage.Ticker, config config.Config) error {
	if !ticker.Twitter.Connected() {
		return errors.New("ticker is not connected to twitter")
	}

	if !config.TwitterEnabled() {
		return errors.New("ticker is not configured with twitter credentials")
	}

	return nil
}
