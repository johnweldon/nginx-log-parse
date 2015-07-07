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
			fmt.Fprintf(os.Stdout, "%20s %16s %d %-40s %s\n", record.TimeLocal.Format("2006-01-02 15:04:05"), record.RemoteAddr, record.Status, record.Request, record.HttpReferrer)
		case log := <-p.Log:
			fmt.Fprintf(os.Stderr, "%s", log)
		case <-p.Dying():
			return
		}
	}
}
