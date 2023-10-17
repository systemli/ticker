package legacy

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mattn/go-mastodon"
	geojson "github.com/paulmach/go.geojson"
)

type Message struct {
	ID             int       `storm:"id,increment"`
	CreationDate   time.Time `storm:"index"`
	Ticker         int       `storm:"index"`
	Text           string
	Attachments    []Attachment
	GeoInformation geojson.FeatureCollection
	Tweet          Tweet
	Telegram       TelegramMeta
	Mastodon       mastodon.Status
}

type Attachment struct {
	UUID        string
	Extension   string
	ContentType string
}

type Tweet struct {
	ID       string
	UserName string
}

type TelegramMeta struct {
	Messages []tgbotapi.Message
}
