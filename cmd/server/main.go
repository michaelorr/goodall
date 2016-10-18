package main

import "log"
import "github.com/michaelorr/goodall/pkg/version"

func main() {
	log.Println(version.VERSION)
	log.Println("Inside the server")
}
