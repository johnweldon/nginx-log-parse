package main

import (
	"fmt"
	"io"
	"os"

	"github.com/johnweldon/nginx-log-parse/nginx"
	"github.com/johnweldon/nginx-log-parse/util"
)

func main() {
	referrers := getReferrers(os.Stdin)

	for r, c := range referrers {
		fmt.Fprintf(os.Stdout, "%-40s :: %d\n", r, c)
	}
}

func getReferrers(in io.Reader) map[string]int {
	p := util.NewParser(in)
	referrers := map[string]int{}
	for {
		select {
		case line := <-p.LineCh:
			if line != nil {
				if rl, ok := line.(nginx.RequestLine); ok {
					referrer := rl.RequestHttpReferrer()
					referrers[referrer] += 1
				}
			}
		case <-p.Dying():
			return referrers
		}
	}
}
