package cli

import (
	"github.com/michaelorr/goodall/pkg/agent"
)

func Run() int {
	c := parseArgs()
	return agent.Run(c.MetricIntervalMs, c.RetentionMin, c.DBPath)
}
