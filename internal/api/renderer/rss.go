package renderer

import (
	"net/http"

	"github.com/gorilla/feeds"
)

const (
	AtomFormat Format = "atom"
	RSSFormat  Format = "rss"
)

type Format string

func FormatFromString(format string) Format {
	if format == string(AtomFormat) {
		return AtomFormat
	}

	return RSSFormat
}

type Feed struct {
	Format Format
	Data   *feeds.Feed
}

var feedContentType = []string{"application/xml; charset=utf-8"}

func (r Feed) Render(w http.ResponseWriter) error {
	return WriteFeed(w, r.Data, r.Format)
}

func (r Feed) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, feedContentType)
}

func WriteFeed(w http.ResponseWriter, data *feeds.Feed, format Format) error {
	var feed string
	var err error

	writeContentType(w, feedContentType)
	if format == AtomFormat {
		feed, err = data.ToAtom()
	} else {
		feed, err = data.ToRss()
	}
	if err != nil {
		log.WithError(err).Error("failed to generate atom")
		return err
	}

	_, err = w.Write([]byte(feed))
	if err != nil {
		log.WithError(err).Error("failed to write response")
		return err
	}

	return nil
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
