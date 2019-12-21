package model_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
)

func TestLoadConfigWithoutPath(t *testing.T) {
	err := os.Setenv("TICKER_LISTEN", ":7070")
	if err != nil {
		t.Fail()
	}

	c := model.LoadConfig("")
	assert.Equal(t, ":7070", c.Listen)

	err = os.Unsetenv("TICKER_LISTEN")
	if err != nil {
		t.Fail()
	}
}

func TestLoadConfigWithPath(t *testing.T) {
	c := model.LoadConfig("../../config.yml")
	assert.Equal(t, ":8080", c.Listen)

	c = model.LoadConfig("config.yml")
	assert.Equal(t, ":8080", c.Listen)
}

func TestConfig_TwitterEnabled(t *testing.T) {
	c := model.NewConfig()

	assert.False(t, c.TwitterEnabled())

	c.TwitterConsumerKey = "a"
	c.TwitterConsumerSecret = "a"

	assert.True(t, c.TwitterEnabled())
}
