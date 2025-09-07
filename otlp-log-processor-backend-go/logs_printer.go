package main

import "fmt"

type stdoutPrinter struct{}

func newStdoutPrinter() *stdoutPrinter {
	return &stdoutPrinter{}
}

func (s stdoutPrinter) print(counts map[any]int64) {
	for k, v := range counts {
		fmt.Printf("%s - %d\n", k, v)
	}
}
