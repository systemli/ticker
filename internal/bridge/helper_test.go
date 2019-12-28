package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/bridge"
	"github.com/systemli/ticker/internal/model"
)

func TestTwitterClient(t *testing.T) {
	setupHelperTestData()
	c := bridge.TwitterClient("a", "b")

	assert.NotNil(t, c)
}

func TestTwitterConnectionEnabled(t *testing.T) {
	model.Config = model.NewConfig()
	ticker := model.NewTicker()

	assert.NotNil(t, bridge.TwitterConnectionEnabled(ticker))

	setupHelperTestData()

	assert.NotNil(t, bridge.TwitterConnectionEnabled(ticker))

	ticker.Twitter.Token = "token"
	ticker.Twitter.Secret = "secret"

	assert.Nil(t, bridge.TwitterConnectionEnabled(ticker))
}

func setupHelperTestData() {
	model.Config = model.NewConfig()
	model.Config.TwitterConsumerSecret = "consumer_secret"
	model.Config.TwitterConsumerKey = "consumer_key"
}
