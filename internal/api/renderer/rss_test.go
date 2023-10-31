package renderer

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/suite"
)

type RendererTestSuite struct {
	suite.Suite
}

func (s *RendererTestSuite) TestFormatFromString() {
	s.Run("when format is atom", func() {
		format := FormatFromString("atom")
		s.Equal(AtomFormat, format)
	})

	s.Run("when format is rss", func() {
		format := FormatFromString("rss")
		s.Equal(RSSFormat, format)
	})

	s.Run("when format is empty", func() {
		format := FormatFromString("")
		s.Equal(RSSFormat, format)
	})
}

func (s *RendererTestSuite) TestWriteFeed() {
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

	s.Run("when format is atom", func() {
		w := httptest.NewRecorder()
		atom := Feed{Data: feed, Format: AtomFormat}

		err := atom.Render(w)
		s.NoError(err)

		atom.WriteContentType(w)
		s.Equal("application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	})

	s.Run("when format is rss", func() {
		w := httptest.NewRecorder()
		rss := Feed{Data: feed, Format: RSSFormat}

		err := rss.Render(w)
		s.NoError(err)

		rss.WriteContentType(w)
		s.Equal("application/xml; charset=utf-8", w.Header().Get("Content-Type"))
	})
}

func TestRendererTestSuite(t *testing.T) {
	suite.Run(t, new(RendererTestSuite))
}
