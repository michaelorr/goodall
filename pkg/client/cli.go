package client

import (
	"github.com/michaelorr/goodall/pkg/db"
)

func Run() int {
	conn, err := db.Open()
	if err != nil {
		// TODO log
		return 1
	}
	err = db.Init()
	if err != nil {
		// TODO log
		return 2
	}

	response := make(chan int)
	go GatherMetrics(response)
	// TODO
	// select response
	// return that value
}

func GatherMetrics(killed chan int) {
	// TODO gather metrics
	// store in bolt
	// sleep for one second
}
