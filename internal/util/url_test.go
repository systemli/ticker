package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractURL(t *testing.T) {
	text := "This is a text with a URL https://example.com"
	urls := ExtractURLs(text)

	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "https://example.com", urls[0])

	text = "This is a text with a URL https://example.com and another URL http://example.org"
	urls = ExtractURLs(text)

	assert.Equal(t, 2, len(urls))
	assert.Equal(t, "https://example.com", urls[0])
	assert.Equal(t, "http://example.org", urls[1])

	text = "This is a text without a URL"
	urls = ExtractURLs(text)

	assert.Equal(t, 0, len(urls))

	text = "This is a text with a URL https://www.systemli.org/en/contact/"
	urls = ExtractURLs(text)

	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "https://www.systemli.org/en/contact/", urls[0])

	text = "This is a text with a URL https://www.systemli.org/en/contact/?key=value"
	urls = ExtractURLs(text)

	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "https://www.systemli.org/en/contact/?key=value", urls[0])
}
