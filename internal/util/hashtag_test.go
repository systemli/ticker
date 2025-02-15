package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractHashtags(t *testing.T) {
	testCases := []struct {
		text     string
		expected []string
	}{
		{
			text:     "Hello #world",
			expected: []string{"#world"},
		},
		{
			text:     "Hello #world #foo",
			expected: []string{"#world", "#foo"},
		},
		{
			text:     "Hello world",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		hashtags := ExtractHashtags(tc.text)
		assert.Equal(t, len(tc.expected), len(hashtags))
		for i, hashtag := range hashtags {
			assert.Equal(t, hashtag, tc.expected[i])
		}
	}
}
