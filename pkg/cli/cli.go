package cli

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MetricIntervalMs time.Duration
	RetentionMin     time.Duration
	DBPath           string
}

func parseArgs() *Config {
	var c *Config = &Config{
		MetricIntervalMs: 1000 * time.Millisecond,
		RetentionMin:     240 * time.Minute,
		DBPath:           "goodall.db",
	}

	if os.Getenv("GOODALL_COLLECTION_MS") != "" {
		coll_ms, err := strconv.Atoi(os.Getenv("GOODALL_COLLECTION_MS"))
		if err != nil {
			log.Println("Unable to parse env var GOODALL_COLLECTION_MS. Falling back to default (1000ms).")
		} else if coll_ms <= 0 {
			log.Println("Cannot use collection interval <=0. Falling back to default (1000ms).")
		} else {
			c.MetricIntervalMs = time.Duration(coll_ms) * time.Millisecond
		}
	}

	if os.Getenv("GOODALL_RETENTION_MIN") != "" {
		retention_min, err := strconv.Atoi(os.Getenv("GOODALL_RETENTION_MIN"))
		if err != nil {
			log.Println("Unable to parse env var GOODALL_RETENTION_MIN. Falling back to default (240m).")
		} else if retention_min <= 0 {
			log.Println("Cannot use retention period <=0. Falling back to default (240m).")
		} else {
			c.RetentionMin = time.Duration(retention_min) * time.Minute
		}
	}

	if os.Getenv("GOODALL_DB_PATH") != "" {
		c.DBPath = os.Getenv("GOODALL_DB_PATH")
	}

	return c
}
