package client

import (
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
	for {
		for bucket, fetch_metric := range metrics.BucketMap {
			// TODO do this fetching in goroutines
			val := fetch_metric()
			// TODO: clean this up
			_ = bucket
			_ = val
			// store in db
		}
		time.Sleep(metrics.Interval)
	}
	response <- 0
}
