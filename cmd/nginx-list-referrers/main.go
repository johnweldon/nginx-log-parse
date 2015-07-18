package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/johnweldon/nginx-log-parse/nginx"
	"github.com/johnweldon/nginx-log-parse/parser"
)

func main() {
	referrers := getReferrers(os.Stdin)

	for r, c := range referrers {
		fmt.Fprintf(os.Stdout, "%5d %s\n", c, r)
	}
}

func getReferrers(in io.Reader) map[string]int {
	p := parser.NewLogFileParser(in)
	referrers := map[string]int{}
	for {
		select {
		case line := <-p.LineCh:
			if line != nil {
				if rl, ok := line.(nginx.RequestLine); ok {
					referrer := rl.RequestHttpReferrer()
					if referrer == "-" {
						continue
					}
					if u, err := url.Parse(referrer); err == nil {
						referrers[u.Scheme+"://"+u.Host] += 1
					} else {
						referrers[referrer] += 1
					}
				}
			}
		case <-p.Dying():
			return referrers
		}
	}
}
