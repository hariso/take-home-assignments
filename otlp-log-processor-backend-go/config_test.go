package main

import (
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestParseConfig(t *testing.T) {
	is := is.New(t)

	originalArgs := os.Args
	os.Args = []string{originalArgs[0], "--countWindow=12s"}
	defer func() {
		os.Args = originalArgs
	}()

	cfg, err := parseConfig()
	is.NoErr(err)
	is.Equal(12*time.Second, cfg.countWindow)
}

// todo more tests are need for all configuration parameters
// e.g. ability to parse those, validation, etc.
