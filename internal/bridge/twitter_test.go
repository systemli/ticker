package bridge_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestSendTweet(t *testing.T) {
	model.Config = model.NewConfig()

	ticker := model.NewTicker()
	message := model.NewMessage()

	// Ticker Twitter is disabled
	assert.Nil(t, bridge.SendTweet(ticker, message))

	// Ticker Twitter is enabled but has no creds
	ticker.Twitter.Active = true
	assert.Nil(t, bridge.SendTweet(ticker, message))

	ticker.Twitter.Token = "token"
	ticker.Twitter.Secret = "secret"
	assert.Nil(t, bridge.SendTweet(ticker, message))

	setupTwitterTestData()
	assert.NotNil(t, bridge.SendTweet(ticker, message))
}

func TestSendTweetWithAttachment(t *testing.T) {
	setupTwitterTestData()

	ticker := model.NewTicker()
	ticker.Twitter.Active = true
	ticker.Twitter.Token = "token"
	ticker.Twitter.Secret = "secret"

	attachment := &model.Attachment{UUID: uuid.New().String()}

	message := model.NewMessage()
	message.Attachments = []model.Attachment{*attachment}

	assert.NotNil(t, bridge.SendTweet(ticker, message))

	upload := model.NewUpload("filename.jpg", "image/jpeg", ticker.ID)
	upload.UUID = attachment.UUID
	_ = storage.DB.Save(upload)

	assert.NotNil(t, bridge.SendTweet(ticker, message))
}

func TestDeleteTweet(t *testing.T) {
	setupTwitterTestData()

	ticker := model.NewTicker()
	message := model.NewMessage()

	assert.Nil(t, bridge.DeleteTweet(ticker, message))

	ticker.Twitter.Active = true
	ticker.Twitter.Token = "token"
	ticker.Twitter.Secret = "secret"

	assert.Nil(t, bridge.DeleteTweet(ticker, message))

	message.Tweet.ID = "foobar"

	assert.NotNil(t, bridge.DeleteTweet(ticker, message))

	message.Tweet.ID = "1"

	assert.NotNil(t, bridge.DeleteTweet(ticker, message))
}

func TestTwitterUser(t *testing.T) {
	setupTwitterTestData()

	ticker := model.NewTicker()
	_, err := bridge.TwitterUser(ticker)
	assert.NotNil(t, err)

	ticker.Twitter.Active = true
	_, err = bridge.TwitterUser(ticker)
	assert.NotNil(t, err)

	ticker.Twitter.Token = "token"
	ticker.Twitter.Secret = "secret"
	_, err = bridge.TwitterUser(ticker)
	assert.NotNil(t, err)

}

func setupTwitterTestData() {
	model.Config = model.NewConfig()
	model.Config.TwitterConsumerSecret = "consumer_secret"
	model.Config.TwitterConsumerKey = "consumer_key"

	if storage.DB == nil {
		storage.DB = storage.OpenDB(fmt.Sprintf("%s/ticker_%d.db", os.TempDir(), time.Now().Nanosecond()))
	}
	_ = storage.DB.Drop("Ticker")
	_ = storage.DB.Drop("Message")
	_ = storage.DB.Drop("Upload")
	_ = storage.DB.Drop("User")
	_ = storage.DB.Drop("Setting")
}
