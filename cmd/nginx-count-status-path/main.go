package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/johnweldon/nginx-log-parse/nginx"
	"github.com/johnweldon/nginx-log-parse/parser"
)

func main() {
	p := parser.NewParser(os.Stdin)
	records := p.GetRecords()

	res := map[string]map[string]map[string]int{}
	for _, record := range records {
		rl, ok := record.(nginx.RequestLine)
		if !ok {
			continue
		}
		status := fmt.Sprintf("%d", rl.ResponseStatus())
		method := rl.RequestMethod()
		u := rl.RequestURI()
		path := u
		pu, err := url.Parse(u)
		if err == nil {
			path = pu.Path
		}
		if _, ok := res[status]; !ok {
			res[status] = map[string]map[string]int{}
		}
		if _, ok := res[status][method]; !ok {
			res[status][method] = map[string]int{}
		}
		res[status][method][path] += 1
	}
	js, err := ToJSON(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "\n%s\n", js)
}

func ToJSON(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}
