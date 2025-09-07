package main

import (
	"fmt"
	"strings"
)

type stdoutPrinter struct{}

func newStdoutPrinter() *stdoutPrinter {
	return &stdoutPrinter{}
}

func (s stdoutPrinter) print(counts map[string]int64) {
	var b strings.Builder

	fmt.Fprint(&b, "attribute value counts:\n")
	for k, v := range counts {
		fmt.Fprintf(&b, "%v - %d\n", k, v)
	}

	fmt.Println(b.String())

}
