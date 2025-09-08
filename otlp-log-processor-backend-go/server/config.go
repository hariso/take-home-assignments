package server

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type config struct {
	listenAddr            string
	maxReceiveMessageSize int
	attributeKey          string
	countWindow           time.Duration
}

var (
	listenAddr = flag.String(
		"listenAddr",
		"localhost:4317",
		"The listen address",
	)
	maxReceiveMessageSize = flag.Int(
		"maxReceiveMessageSize",
		16777216,
		"The max message size in bytes the server can receive",
	)
	attributeKey = flag.String(
		"attributeKey",
		"",
		"Attribute key for which the numbers of distinct values should be tracked",
	)
	countWindow = flag.Duration(
		"countWindow",
		10*time.Second,
		"Duration of the time window after which the number of distinct values of attributeKey will be printed.",
	)
)

// parseConfig parses the app's config and validates the values.
// todo add validation for all config parameters
func parseConfig() (config, error) {
	// NB: parseConfig returns a configuration struct,
	// without the caller having to know how exactly it's parsed.
	// We can add parsing the config values from files or env. variables easily.

	flag.Parse()

	cfg := config{
		listenAddr:            *listenAddr,
		maxReceiveMessageSize: *maxReceiveMessageSize,
		attributeKey:          *attributeKey,
		countWindow:           *countWindow,
	}

	err := validateConfig(cfg)
	if err != nil {
		return config{}, fmt.Errorf("validation error: %w", err)
	}

	return cfg, nil
}

func validateConfig(cfg config) error {
	attrKey := strings.TrimSpace(cfg.attributeKey)
	if attrKey == "" {
		return fmt.Errorf("attributeKey is required")
	}

	return nil
}
