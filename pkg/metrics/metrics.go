package metrics

import (
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const Interval time.Duration = 1 * time.Millisecond

type DataPoint struct {
	BucketName string
	Value      float64
}

type metric_f func(string, chan *DataPoint, chan error)

var BucketMap map[string]metric_f = map[string]metric_f{
	"disk_used":      disk_used,
	"disk_free":      disk_free,
	"disk_total":     disk_total,
	"mem_used":       mem_used,
	"mem_available":  mem_available,
	"mem_total":      mem_total,
	"system_load_1":  load1,
	"system_load_15": load15,
	"system_load_5":  load5,
}

//////////
// The repeated code below is begging to be ripped into a single helper func
/////////

func disk_used(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Used)}
	}
}

func disk_free(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Free)}
	}
}

func disk_total(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := disk.Usage("/"); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Total)}
	}
}

func mem_used(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := mem.VirtualMemory(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Used)}
	}
}

func mem_available(bucket string, result chan *DataPoint, errors chan error) {
	if stat, err := mem.VirtualMemory(); err != nil {
		errors <- err
	} else {
		result <- &DataPoint{bucket, float64(stat.Available)}
	}
}

func mem_total(bucket string, result chan *DataPoint, errors chan error) {
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
