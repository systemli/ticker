package bridge

import (
	"errors"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/systemli/ticker/internal/model"
)

//TwitterClient returns a client for twitter api
func TwitterClient(t, s string) *twitter.Client {
	config := oauth1.NewConfig(model.Config.TwitterConsumerKey, model.Config.TwitterConsumerSecret)
	token := oauth1.NewToken(t, s)

	return twitter.NewClient(config.Client(oauth1.NoContext, token))
}

//TwitterConnectionEnabled returns true when ticker can use twitter
func TwitterConnectionEnabled(ticker *model.Ticker) error {
	if !ticker.Twitter.Connected() {
		return errors.New("ticker is not connected to twitter")
	}

	if !model.Config.TwitterEnabled() {
		return errors.New("ticker is not configured with twitter credentials")
	}

	return nil
}
