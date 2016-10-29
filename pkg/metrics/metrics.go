package metrics

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type DataPoint struct {
	Name  string
	Value float64
}

type JsonMetric struct {
	DataPoint
	Timestamp string
}

type JsonPayload struct {
	Timestamp string
	Metrics   []JsonMetric
}

type metricF func(string, chan *DataPoint, chan error)

var BucketMap map[string]metricF = map[string]metricF{
	"disk_used":      diskUsed,
	"disk_free":      diskFree,
	"disk_total":     diskTotal,
	"mem_used":       memUsed,
	"mem_available":  memAvailable,
	"mem_total":      memTotal,
	"system_load_1":  load1,
	"system_load_15": load15,
	"system_load_5":  load5,
}

//////////
// The repeated code below is begging to be ripped into a single helper func
/////////

func diskUsed(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Used)}
	}
}

func diskFree(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Free)}
	}
}

func diskTotal(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Total)}
	}
}

func memUsed(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := mem.VirtualMemory(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Used)}
	}
}

func memAvailable(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := mem.VirtualMemory(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Available)}
	}
}

func memTotal(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := mem.VirtualMemory(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Total)}
	}
}

func load1(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := load.Avg(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, stat.Load15}
	}
}

func load15(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := load.Avg(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, stat.Load1}
	}
}

func load5(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := load.Avg(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, stat.Load5}
	}
}
