package db

import (
	"encoding/json"
	"fmt"
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
				log.Printf("Problem fetching from bucket %s\n", metricName)
				continue
			}

			val_f, err := Btof(val_b)
			if err != nil {
				log.Printf("There was an error converting %s to float64\n", val_b)
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

func BucketPayload(conn *bolt.DB, bucketName string) ([]byte, error) {
	// This func is hardcoded to 10 objects in the response, it would be nice
	// to parameterize this value but that would require using something other
	// than DefaultServeMux
	metricSlice := make([]metrics.JsonMetric, 0)

	if err := conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Unable to fetch bucket: %s", bucketName)
		}

		c := b.Cursor()
		k, v := c.Last()
		metricSlice = insertIntoPayload(metricSlice, bucketName, k, v)

		for i := 0; i < 9; i++ {
			k, v := c.Prev()
			metricSlice = insertIntoPayload(metricSlice, bucketName, k, v)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	response := metrics.JsonPayload{time.Now().UTC().Format("2006-01-02T15:04:05.999"), metricSlice}
	return json.Marshal(response)
}

func insertIntoPayload(metricSlice []metrics.JsonMetric, metricName string, key_b, val_b []byte) []metrics.JsonMetric {
	// if Last() or Prev() don't return valid data, we can't insert it
	if key_b == nil {
		log.Println("Error fetching value, skipping")
		return metricSlice
	}

	val_f, err := Btof(val_b)
	if err != nil {
		log.Printf("There was an error converting %s to float64\n", val_b)
		return metricSlice
	}

	data := metrics.JsonMetric{
		DataPoint: metrics.DataPoint{
			Name:  metricName,
			Value: val_f,
		},
		Timestamp: string(key_b),
	}
	return append(metricSlice, data)
}
