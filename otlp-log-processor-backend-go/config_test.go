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

	cfg := parseConfig()
	is.Equal(12*time.Second, cfg.countWindow)
}
