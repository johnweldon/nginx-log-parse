package main

import (
	"fmt"
	"os"

	"github.com/johnweldon/nginx-log-parse/parser"
)

func main() {
	p := parser.NewEngine(os.Stdin)
	for {
		select {
		case line := <-p.LogLines():
			if line != nil {
				fmt.Fprintf(os.Stdout, "%s\n", line)
			}
		case <-p.Dying():
			return
		}
	}
}
