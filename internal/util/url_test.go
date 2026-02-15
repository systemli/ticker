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
		{
			"Text with URL in parentheses (https://www.systemli.org) and more text",
			[]string{"https://www.systemli.org"},
		},
		{
			"Wikipedia URL https://en.wikipedia.org/wiki/Foo_(bar) should keep balanced parens",
			[]string{"https://en.wikipedia.org/wiki/Foo_(bar)"},
		},
		{
			"Wrapped Wikipedia URL (https://en.wikipedia.org/wiki/Foo_(bar)) should trim outer paren",
			[]string{"https://en.wikipedia.org/wiki/Foo_(bar)"},
		},
		{
			"URL with trailing dot https://example.com.",
			[]string{"https://example.com"},
		},
		{
			"URL with trailing comma https://example.com, and more",
			[]string{"https://example.com"},
		},
		{
			"URL with trailing exclamation https://example.com!",
			[]string{"https://example.com"},
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
