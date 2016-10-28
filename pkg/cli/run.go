package cli

import (
	"log"

	"github.com/michaelorr/goodall/pkg/agent"
	"github.com/michaelorr/goodall/pkg/db"
	"github.com/michaelorr/goodall/pkg/server"
)

func Run() int {
	c := parseArgs()

	conn, err := db.Open(c.DBPath)
	if err != nil {
		log.Println(err)
		return 1
	}
	err = db.Init(conn)
	if err != nil {
		log.Println(err)
		return 2
	}

	ret_val := make(chan int)
	go agent.Run(conn, c.MetricIntervalMs, c.RetentionMin)
	go server.Run(conn, c.HTTPPort, ret_val)

	return <-ret_val
}
