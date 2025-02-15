package util

import "strings"

func ExtractHashtags(text string) []string {
	var hashtags []string

	for _, word := range strings.Fields(text) {
		if strings.HasPrefix(word, "#") {
			hashtags = append(hashtags, word)
		}
	}

	return hashtags
}
