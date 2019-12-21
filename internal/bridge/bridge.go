package bridge

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/systemli/ticker/internal/model"
	"strconv"
)

var Twitter *TwitterBridge

//
type TwitterBridge struct {
	ConsumerKey    string
	ConsumerSecret string
}

//
func NewTwitterBridge(key, secret string) *TwitterBridge {
	return &TwitterBridge{
		ConsumerKey:    key,
		ConsumerSecret: secret,
	}
}

//
func (tb *TwitterBridge) Update(ticker model.Ticker, message model.Message) (*twitter.Tweet, error) {
	client := tb.client(ticker.Twitter.Token, ticker.Twitter.Secret)

	tweet, _, err := client.Statuses.Update(message.PrepareTweet(&ticker), nil)
	if err != nil {
		return tweet, err
	}

	return tweet, nil
}

func (tb *TwitterBridge) Delete(ticker model.Ticker, tweetID string) error {
	client := tb.client(ticker.Twitter.Token, ticker.Twitter.Secret)

	id, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		return err
	}

	_, _, err = client.Statuses.Destroy(id, nil)

	return err
}

//User returns the user information.
func (tb *TwitterBridge) User(ticker model.Ticker) (*twitter.User, error) {
	token := oauth1.NewToken(ticker.Twitter.Token, ticker.Twitter.Secret)
	httpClient := tb.config().Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		IncludeEmail:    twitter.Bool(false),
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
	})

	if err != nil {
		return user, err
	}

	return user, nil
}

func (tb *TwitterBridge) config() *oauth1.Config {
	return oauth1.NewConfig(tb.ConsumerKey, tb.ConsumerSecret)
}

func (tb *TwitterBridge) client(accessToken, accessSecret string) *twitter.Client {
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := tb.config().Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}
