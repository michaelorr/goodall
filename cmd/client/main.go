package main

import (
	"os"

	"github.com/michaelorr/goodall/pkg/client"
)

func main() {
	os.Exit(client.Run())
}
