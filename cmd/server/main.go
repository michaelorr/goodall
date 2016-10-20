package main

import (
	"os"

	"github.com/michaelorr/goodall/pkg/server"
)

func main() {
	os.Exit(server.Run())
}
