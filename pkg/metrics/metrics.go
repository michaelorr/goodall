package metrics

import "time"

const Interval time.Duration = 1 * time.Millisecond

type DataPoint struct {
	BucketName string
	Value      float64
}

type metric_f func(string, chan *DataPoint)

var BucketMap map[string]metric_f = map[string]metric_f{
	"cpu_usage":       cpu,
	"io_wait":         iowait,
	"memory":          memory,
	"network_traffic": network,
	"system_load_1":   load1,
	"system_load_15":  load15,
	"system_load_5":   load5,
}

func cpu(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 0}
}

func iowait(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 1}
}

func memory(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 2}
}

func network(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 3}
}

func load1(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 4}
}

func load15(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 5}
}

func load5(bucket string, result chan *DataPoint) {
	result <- &DataPoint{bucket, 6}
}
