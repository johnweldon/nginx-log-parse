package main

import (
	"fmt"
	"os"
	"time"
)

type LogLine struct {
	RemoteAddr    string
	RemoteUser    string
	TimeLocal     time.Time
	Request       string
	Status        int
	BodyBytesSent int
	HttpReferrer  string
	HttpUserAgent string
}

func main() {
	p := NewParser(os.Stdin)
	records := p.GetRecords()
	fmt.Fprintf(os.Stdout, "%#v\n", records)
}
