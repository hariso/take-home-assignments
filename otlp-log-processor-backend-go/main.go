package main

import (
	"fmt"
	"log"

	"dash0.com/otlp-log-processor-backend/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to run server: %w", err))
	}
}
