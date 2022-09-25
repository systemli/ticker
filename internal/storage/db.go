package storage

import "github.com/asdine/storm"

//OpenDB returns a storm.DB reference
func OpenDB(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		log.WithError(err).Panic("failed to open database file")
	}

	return db
}
