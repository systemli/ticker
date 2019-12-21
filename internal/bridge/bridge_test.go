package bridge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/bridge"
)

func TestNewTwitterBridge(t *testing.T) {
	tb := bridge.NewTwitterBridge("key", "secret")

	assert.Equal(t, "key", tb.ConsumerKey)
	assert.Equal(t, "secret", tb.ConsumerSecret)
}
