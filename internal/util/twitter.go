package util

import (
	"fmt"

	"github.com/systemli/ticker/internal/model"
)

//PrepareTweet builds the string for the Tweet.
func PrepareTweet(ticker *model.Ticker, message *model.Message) string {
	tweet := message.Text
	if ticker.PrependTime {
		tweet = fmt.Sprintf(`%.2d:%.2d %s`, message.CreationDate.Hour(), message.CreationDate.Minute(), tweet)
	}

	//TODO: Add default hashtags
	//TODO: Check length, split long tweets

	return tweet
}
