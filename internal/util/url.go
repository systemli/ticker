package util

import "regexp"

// ExtractURLs extracts URLs from a text.
func ExtractURLs(text string) []string {
	urlRegex := regexp.MustCompile(`(http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+)`)
	return urlRegex.FindAllString(text, -1)
}
