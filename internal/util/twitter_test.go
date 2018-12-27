package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/util"
)

func TestPrepareTweet(t *testing.T) {
	ticker := model.NewTicker()
	message := model.NewMessage()
	message.CreationDate, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	message.Text = "example"

	assert.Equal(t, "example", util.PrepareTweet(ticker, message))

	ticker.PrependTime = true

	assert.Equal(t, "22:08 example", util.PrepareTweet(ticker, message))
}
