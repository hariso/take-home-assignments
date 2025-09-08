package main

import (
	"context"
	"fmt"
	"log"

	"dash0.com/otlp-log-processor-backend/server"
)

func main() {
	err := server.Run(context.Background())
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to run server: %w", err))
	}
}
