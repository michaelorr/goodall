package main

import (
	"os"

	"github.com/michaelorr/goodall/pkg/agent"
)

func main() {
	os.Exit(agent.Run())
}
