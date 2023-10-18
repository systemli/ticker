package legacy

import "time"

type Upload struct {
	ID           int       `storm:"id,increment"`
	UUID         string    `storm:"index,unique"`
	CreationDate time.Time `storm:"index"`
	TickerID     int       `storm:"index"`
	Path         string
	Extension    string
	ContentType  string
}
