package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithoutPath(t *testing.T) {
	err := os.Setenv("TICKER_LISTEN", ":7070")
	if err != nil {
		t.Fail()
	}

	c := LoadConfig("")
	assert.Equal(t, ":7070", c.Listen)

	err = os.Unsetenv("TICKER_LISTEN")
	if err != nil {
		t.Fail()
	}
}

func TestLoadConfigWithPath(t *testing.T) {
	c := LoadConfig("../../config.yml")
	assert.Equal(t, ":8080", c.Listen)

	c = LoadConfig("config.yml")
	assert.Equal(t, ":8080", c.Listen)
}

func TestLoadConfigWithFallback(t *testing.T) {
	c := LoadConfig("/x/y/z")
	assert.Equal(t, ":8080", c.Listen)
}

func TestConfig_TwitterEnabled(t *testing.T) {
	c := NewConfig()

	assert.False(t, c.TwitterEnabled())

	c.TwitterConsumerKey = "a"
	c.TwitterConsumerSecret = "a"

	assert.True(t, c.TwitterEnabled())
}

func TestConfig_TelegramEnabled(t *testing.T) {
	c := NewConfig()

	assert.False(t, c.TelegramEnabled())

	c.TelegramBotToken = "a"

	assert.True(t, c.TelegramEnabled())
}
