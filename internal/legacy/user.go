package legacy

import "time"

type User struct {
	ID                int       `storm:"id,increment"`
	CreationDate      time.Time `storm:"index"`
	Email             string    `storm:"unique"`
	Role              string
	EncryptedPassword string
	IsSuperAdmin      bool
	Tickers           []int
}
