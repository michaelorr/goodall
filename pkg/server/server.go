package server

import (
	"fmt"
	"log"
	"net/http"
)

func Run(port int, ret_val chan int) {
	http.HandleFunc("/", hello)

	log.Printf("listening on %d. Go to http://127.0.0.1:%d/\n", port, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}
