package renderer

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	w := httptest.NewRecorder()

	feed := &feeds.Feed{
		Title: "Title",
		Author: &feeds.Author{
			Name:  "Name",
			Email: "Email",
		},
		Link: &feeds.Link{
			Href: "https://demoticker.org",
		},
		Created: time.Now(),
	}
	atom := Feed{Data: feed, Format: AtomFormat}

	err := atom.Render(w)
	assert.Nil(t, err)

	atom.WriteContentType(w)
	assert.Equal(t, "application/xml; charset=utf-8", w.Header().Get("Content-Type"))

	rss := Feed{Data: feed, Format: RSSFormat}

	err = rss.Render(w)
	assert.Nil(t, err)
}

func TestFormatFromString(t *testing.T) {
	var format Format

	format = FormatFromString("atom")
	assert.Equal(t, AtomFormat, format)

	format = FormatFromString("rss")
	assert.Equal(t, RSSFormat, format)

	format = FormatFromString("")
	assert.Equal(t, RSSFormat, format)
}
