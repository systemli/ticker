package storage

import "github.com/asdine/storm"

var (
	DB *storm.DB
)

//OpenDB returns a storm.DB reference
func OpenDB(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		panic(err)
	}

	return db
}
