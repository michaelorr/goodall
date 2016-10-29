package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func LatestPayload(conn *bolt.DB) ([]byte, error) {
	metricSlice := make([]metrics.JsonMetric, 0)

	conn.View(func(tx *bolt.Tx) error {
		for metricName, _ := range metrics.BucketMap {
			key_b, val_b := LatestFromBucket(tx, metricName)
			if key_b == nil {
				log.Println("Problem fetching from bucket %s", metricName)
				continue
			}

			val_f, err := Btof(val_b)
			if err != nil {
				log.Printf("There was an error converting %s to float64", val_b)
				continue
			}

			data := metrics.JsonMetric{
				DataPoint: metrics.DataPoint{
					Name:  metricName,
					Value: val_f,
				},
				Timestamp: string(key_b),
			}
			metricSlice = append(metricSlice, data)
		}
		return nil
	})

	response := metrics.JsonPayload{time.Now().UTC().Format("2006-01-02T15:04:05.999"), metricSlice}
	return json.Marshal(response)
}

func LatestFromBucket(tx *bolt.Tx, bucketName string) ([]byte, []byte) {
	b := tx.Bucket([]byte(bucketName))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	return c.Last()
}
