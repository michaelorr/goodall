package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func Store(conn *bolt.DB, result *metrics.DataPoint, now string, wg *sync.WaitGroup) {
	defer wg.Done()

	if err := conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(result.Name))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", result.Name)
		}

        val, err := Ftob(result.Value)
        if err != nil {
            return err
        }
        return b.Put([]byte(now), val)
	}); err != nil {
		log.Println(err)
	}
}
