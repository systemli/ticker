package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractURL(t *testing.T) {
	testCases := []struct {
		text     string
		expected []string
	}{
		{
			"This is a text with a URL https://example.com",
			[]string{"https://example.com"},
		},
		{
			"This is a text with a URL https://example.com and another URL http://example.org",
			[]string{"https://example.com", "http://example.org"},
		},
		{
			"This is a text without a URL",
			[]string{},
		},
		{
			"This is a text with a URL https://www.systemli.org/en/contact/",
			[]string{"https://www.systemli.org/en/contact/"},
		},
		{
			"This is a text with a URL https://www.systemli.org/en/contact/?key=value",
			[]string{"https://www.systemli.org/en/contact/?key=value"},
		},
		{
			"This is a text with a URL https://www.systemli.org/en/contact/?key=value#fragment",
			[]string{"https://www.systemli.org/en/contact/?key=value#fragment"},
		},
	}

	for _, tc := range testCases {
		urls := ExtractURLs(tc.text)
		assert.Equal(t, len(tc.expected), len(urls))
		for i, url := range tc.expected {
			assert.Equal(t, url, urls[i])
		}
	}
}
