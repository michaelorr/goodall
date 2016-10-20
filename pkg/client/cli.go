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
	err = db.Init(conn)
	if err != nil {
		// TODO log
		return 2
	}

	response := make(chan int)
	go GatherMetrics(response)
	// TODO
	// select response
	// return that value
	return 0
}

func GatherMetrics(killed chan int) {
	// TODO gather metrics
	// store in bolt
	// sleep for one second
}
