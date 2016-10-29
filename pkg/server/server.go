package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/db"
	"github.com/michaelorr/goodall/pkg/metrics"
)

func Run(conn *bolt.DB, port int) {
	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		response, err := db.LatestPayload(conn)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", response)
	})

	for bucket, _ := range metrics.BucketMap {
		bucketName := bucket // needed to bring this arg into the closure
		http.HandleFunc(fmt.Sprintf("/%s", bucket), func(w http.ResponseWriter, r *http.Request) {
			response, err := db.BucketPayload(conn, bucketName)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%s", response)
		})
	}

	log.Printf("listening on http://127.0.0.1:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
