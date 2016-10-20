package metrics

import "time"

const Interval time.Duration = 1 * time.Second

var BucketMap map[string]func() float64 = map[string]func() float64{
	"cpu_usage":       CPU,
	"io_wait":         IOwait,
	"memory":          Memory,
	"network_traffic": Network,
	"system_load_1":   Load1,
	"system_load_15":  Load15,
	"system_load_5":   Load5,
}

func CPU() float64 {
	return 0
}

func IOwait() float64 {
	return 1
}

func Memory() float64 {
	return 2
}

func Network() float64 {
	return 3
}

func Load1() float64 {
	return 4
}

func Load15() float64 {
	return 5
}

func Load5() float64 {
	return 6
}
