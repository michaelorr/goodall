package db

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func LatestPayload(conn *bolt.DB) []byte {
	metricSlice := make([]metrics.JsonMetric, 0)

	conn.View(func(tx *bolt.Tx) error {
		for metricName, _ := range metrics.BucketMap {
			key_b, val_b := LatestFromBucket(tx, metricName)
			// TODO handle errors
			val_f, _ := Btof(val_b)

			data := metrics.JsonMetric{
				DataPoint: metrics.DataPoint{
					Name: metricName,
					Value: val_f,
				},
				Timestamp: string(key_b),
			}
			metricSlice = append(metricSlice, data)
		}
		return nil
	})

	response := metrics.JsonPayload{time.Now().String(), metricSlice}
	r, _ := json.Marshal(response)
	// TODO error checking
	return r
}

func LatestFromBucket(tx *bolt.Tx, bucketName string) ([]byte, []byte) {
	b := tx.Bucket([]byte(bucketName))
	// TODO error and nil handling here
	c := b.Cursor()
	return c.Last()
}
