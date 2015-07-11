package main

import (
	"fmt"
	"os"

	"github.com/johnweldon/nginx-log-parse/nginx"
	"github.com/johnweldon/nginx-log-parse/util"
)

func main() {
	p := util.NewParser(os.Stdin)
	for {
		select {
		case line := <-p.LineCh:
			record, ok := line.(nginx.RequestLine)
			if !ok {
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", record)
		case <-p.Dying():
			return
		}
	}
}
