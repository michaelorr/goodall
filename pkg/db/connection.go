package db

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func Open(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
}

func Init(conn *bolt.DB) error {
	for bucket, _ := range metrics.BucketMap {
		err := conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func Ftob(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

func Btof(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	float := math.Float64frombits(bits)
	return float
}
