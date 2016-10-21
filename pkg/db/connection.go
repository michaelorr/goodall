package db

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"

	"github.com/michaelorr/goodall/pkg/metrics"
)

func Open() (*bolt.DB, error) {
	return bolt.Open("goodall.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
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

func Ftob(f float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Btof(b []byte) (float64, error) {
	var f float64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}
