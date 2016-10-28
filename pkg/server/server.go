package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/db"
)

func Run(conn *bolt.DB, port int, ret_val chan int) {
	// TODO make sure we need ret_val

	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving a request on /latest")
		fmt.Fprintf(w, "%s", db.LatestPayload(conn))
	})

	log.Printf("listening on http://127.0.0.1:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
