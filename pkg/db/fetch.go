package db

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func LatestPayload(conn *bolt.DB) string {
	var response string

	conn.View(func(tx *bolt.Tx) error {
		for metricName, _ := range metrics.BucketMap {
			key_b, val_b := LatestFromBucket(tx, metricName)
			// TODO handle errors
			val_f, _ := Btof(val_b)

			// TODO create a struct a Marshall to json
			response = fmt.Sprintf("%s\n%s\t%f\t%s", response, metricName, val_f, string(key_b))
		}

		return nil
	})

	return response
}

func LatestFromBucket(tx *bolt.Tx, bucketName string) ([]byte, []byte) {
	b := tx.Bucket([]byte(bucketName))
	// TODO error and nil handling here
	c := b.Cursor()
	return c.Last()
}
