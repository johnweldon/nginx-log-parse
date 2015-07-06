package main

import (
	"fmt"
	"os"

	"github.com/johnweldon/nginx-log-parse/util"
)

func main() {
	p := util.NewParser(os.Stdin)
	records := p.GetRecords()
	fmt.Fprintf(os.Stdout, "%#v\n", records)
}
