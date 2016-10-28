package cli

import (
	"github.com/michaelorr/goodall/pkg/agent"
	"github.com/michaelorr/goodall/pkg/server"
)

func Run() int {
	c := parseArgs()

	ret_val := make(chan int)
	go agent.Run(c.MetricIntervalMs, c.RetentionMin, c.DBPath, ret_val)
	go server.Run(c.HTTPPort, ret_val)

	return <-ret_val
}
