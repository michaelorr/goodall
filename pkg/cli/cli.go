package cli

import "time"

type Config struct {
	MetricIntervalMs time.Duration
	RetentionMin     time.Duration
	DBPath           string
}

func parseArgs() *Config {
	// TODO pull from command line to parametrize
	return &Config{
		MetricIntervalMs: 1000 * time.Millisecond,
		RetentionMin:     240 * time.Minute,
		DBPath:           "goodall.db",
	}
}
