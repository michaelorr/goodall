package agent

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/db"
	"github.com/michaelorr/goodall/pkg/metrics"
)

var cleanupKeyMin = []byte("2016-01-01T00:00:00Z")

func Run(metricInterval, retentionPeriod time.Duration, path string, ret_val chan int) {
	conn, err := db.Open(path)
	if err != nil {
		log.Println(err)
		ret_val <- 1
		return
	}
	err = db.Init(conn)
	if err != nil {
		log.Println(err)
		ret_val <- 2
		return
	}

	go GatherMetrics(conn, ret_val, metricInterval)
	go CleanupMetrics(conn, metricInterval, retentionPeriod)
}

func CleanupMetrics(conn *bolt.DB, metricInterval, retentionPeriod time.Duration) {
	for {
		cleanupKeyMax := []byte(time.Now().UTC().Add(-retentionPeriod).Format("2006-01-02T15:04:05.999"))

		conn.Update(func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				c := b.Cursor()

				for k, _ := c.Seek(cleanupKeyMin); k != nil && bytes.Compare(k, cleanupKeyMax) <= 0; k, _ = c.Next() {
					err := b.Delete(k)
					if err != nil {
						return err
					}
				}
				return nil
			})
		})

		time.Sleep(metricInterval)
	}
}

func GatherMetrics(conn *bolt.DB, response chan int, metricInterval time.Duration) {
	for {
		var wg sync.WaitGroup
		now := time.Now().UTC().Format("2006-01-02T15:04:05.999")
		results := make(chan *metrics.DataPoint, len(metrics.BucketMap))
		errors := make(chan error)

		// spin off goroutines to fetch each metric
		for bucket, fetchMetric := range metrics.BucketMap {
			wg.Add(1)
			go fetchMetric(bucket, results, errors)
		}

		// TODO handle errors from metrics gathering

		// wait until all metrics goroutines complete before continuing
		go func() {
			wg.Wait()
			close(results)
		}()

		// gather the results
		for result := range results {
			// TODO do this in a separate goroutine in the connection package
			// TODO and defer the call to wg.Done() to make sure we don't leak goroutines
			if err := conn.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(result.BucketName))
				if b == nil {
					return fmt.Errorf("bucket %s does not exist", result.BucketName)
				}
				val, err := db.Ftob(result.Value)
				if err != nil {
					return err
				}
				return b.Put([]byte(now), val)
			}); err != nil {
				log.Println(err)
			}

			wg.Done()
		}

		time.Sleep(metricInterval)
	}
}
