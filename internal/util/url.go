package util

import (
	"regexp"
	"strings"
)

// ExtractURLs extracts URLs from a text.
func ExtractURLs(text string) []string {
	urlRegex := regexp.MustCompile(`(https?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(),]|%[0-9a-fA-F][0-9a-fA-F]|#)+)`)
	matches := urlRegex.FindAllString(text, -1)

	var urls []string
	for _, match := range matches {
		urls = append(urls, cleanURL(match))
	}

	return urls
}

// cleanURL removes trailing punctuation and unbalanced closing brackets from a URL.
func cleanURL(url string) string {
	// Iteratively strip trailing characters that are likely not part of the URL.
	for len(url) > 0 {
		last := url[len(url)-1]

		// Remove trailing punctuation that is almost never part of a URL at the end.
		if last == '.' || last == ',' || last == '!' || last == ';' || last == ':' {
			url = url[:len(url)-1]
			continue
		}

		// Remove trailing closing brackets only if they are unbalanced within the URL path
		// (i.e. there is no matching opening bracket). This preserves URLs like
		// https://en.wikipedia.org/wiki/Foo_(bar) while still trimming the outer bracket
		// in text like (https://example.com).
		if last == ')' {
			path := url[strings.Index(url, "://")+3:]
			if strings.Count(path, "(") < strings.Count(path, ")") {
				url = url[:len(url)-1]
				continue
			}
		}
		if last == ']' {
			path := url[strings.Index(url, "://")+3:]
			if strings.Count(path, "[") < strings.Count(path, "]") {
				url = url[:len(url)-1]
				continue
			}
		}

		break
	}

	return url
}
