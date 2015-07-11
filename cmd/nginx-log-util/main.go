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
		case line := <-p.LineCh:
			if line != nil {
				fmt.Fprintf(os.Stdout, "%s\n", line)
			}
		case <-p.Dying():
			return
		}
	}
}
