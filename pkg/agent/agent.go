package agent

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/db"
	"github.com/michaelorr/goodall/pkg/metrics"
)

// Would be more intelligent to somehow determine the timestamp of the earliest
// entry rather than hardcoding something that is assumed to be in the past
var cleanupKeyMin = []byte("2016-01-01T00:00:00Z")

func Run(conn *bolt.DB, metricInterval, retentionPeriod time.Duration) {
	go GatherMetrics(conn, metricInterval)
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

func GatherMetrics(conn *bolt.DB, metricInterval time.Duration) {
	errors := make(chan error)

	// log any errors from the metrics gathering
	go func() {
		for err := range errors {
			log.Println(err)
		}
	}()

	var wg sync.WaitGroup
	for {
		now := time.Now().UTC().Format("2006-01-02T15:04:05.999")
		results := make(chan *metrics.DataPoint, len(metrics.BucketMap))

		// spin off goroutines to fetch each metric
		for bucket, fetchMetric := range metrics.BucketMap {
			wg.Add(1)
			go fetchMetric(bucket, results, errors)
		}

		// wait until all metrics goroutines complete before continuing
		go func() {
			wg.Wait()
			close(results)
		}()

		// gather and store the results
		for result := range results {
			db.Store(conn, result, now, &wg)
		}

		time.Sleep(metricInterval)
	}
}
