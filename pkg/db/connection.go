package db

import (
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func Open() (*bolt.DB, error) {
	return bolt.Open("goodall.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
}

func Init(conn *bolt.DB) error {
	for bucket, _ := range metrics.BucketMap {
		err := conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}
