package db

import (
	"time"

	"github.com/boltdb/bolt"
)

var BucketNames []string = []string{
	"",
}

func Open() (*bolt.DB, error) {
	return bolt.Open("goodall.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
}

func Init(*bolt.DB) error {
	// make sure that the buckets exist
	return nil
}
