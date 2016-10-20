package client

import (
	"sync"
	"time"

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
	go GatherMetrics(response)
	return <-response
}

func GatherMetrics(response chan int) {
	var results chan *metrics.DataPoint
	var wg sync.WaitGroup

	for {
		results = make(chan *metrics.DataPoint, len(metrics.BucketMap))

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
			// TODO write the result to the DB
			_ = result
			wg.Done()
		}

		time.Sleep(metrics.Interval)
	}
}
