package main

import (
	"fmt"
	"os"

	"github.com/johnweldon/nginx-log-parse/util"
)

func main() {
	p := util.NewParser(os.Stdin)
	for {
		select {
		case record := <-p.EntryCh:
			if record == nil {
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", record)
		case log := <-p.Log:
			fmt.Fprintf(os.Stderr, "%s\n", log)
		case <-p.Dying():
			return
		}
	}
}
