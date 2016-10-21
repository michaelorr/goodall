package client

import (
	"bytes"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/db"
	"github.com/michaelorr/goodall/pkg/metrics"
)

func Run() int {
	conn, err := db.Open()
	if err != nil {
		return 1
	}
	err = db.Init(conn)
	if err != nil {
		return 2
	}

	response := make(chan int)
	go GatherMetrics(conn, response)
	go CleanupMetrics(conn)
	return <-response
}

func CleanupMetrics(conn *bolt.DB) {
	for {
		conn.Update(func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bolt.Bucket) error{
				c := b.Cursor()

				// TODO
				// Extract this to a better location
				min := []byte("2016-01-01T00:00:00Z")
				// TODO
				// Extract this to a better location
				// Make this configurable
				max := []byte(time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339))

				for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
					err := b.Delete(k)
					if err != nil {
						return err
					}
				}
				return nil
			})
		})

		time.Sleep(metrics.Interval)
	}
}

func GatherMetrics(conn *bolt.DB, response chan int) {
	for {
		var wg sync.WaitGroup
		now := time.Now().UTC().Format(time.RFC3339)
		results := make(chan *metrics.DataPoint, len(metrics.BucketMap))

		// spin off goroutines to fetch each metric
		for bucket, fetch_metric := range metrics.BucketMap {
			wg.Add(1)
			go fetch_metric(bucket, results)
		}

		// wait until all metrics goroutines complete before continuing
		go func() {
			wg.Wait()
			close(results)
		}()

		// gather the results
		for result := range results {
			// TODO do this in a separate goroutine in the connection package
			err := conn.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(result.BucketName))
				if b == nil {
					// TODO Bucket does not exist
				}
				val, err := db.Ftob(result.Value)
				// TODO error checking
				err = b.Put([]byte(now), val)
				return err
			})
			// TODO error checking
			_ = err

			wg.Done()
		}

		time.Sleep(metrics.Interval)
	}
}
