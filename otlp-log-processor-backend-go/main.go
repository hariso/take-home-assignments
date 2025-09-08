package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		printHelp()
		os.Exit(1)
	}

	err = run(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func printHelp() {
	// todo pretty print the usage instructions
	fmt.Println("Usage: TODO")
}
