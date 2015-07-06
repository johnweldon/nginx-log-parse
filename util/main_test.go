package util_test

import (
	"bytes"
	"testing"

	. "github.com/johnweldon/nginx-log-parse/util"
)

func TestParse(t *testing.T) {
	p := NewParser(bytes.NewReader([]byte(goodLog)))

	lines := p.GetRecords()
	if len(lines) != 45 {
		t.Errorf("Missing log lines, expected 45 but got %d:\n%#v\n", len(lines), lines)
	}

	if lines[0].RemoteAddr != "10.10.10.15" {
		t.Errorf("Unexpected RemoteAddr, %q", lines[0].RemoteAddr)
	}

	if lines[7].HttpUserAgent[13:22] != "Macintosh" {
		t.Errorf("Unexpected User Agent, %q", lines[7].HttpUserAgent[13:22])
	}

	if lines[13].Request != "GET /filter/tips HTTP/1.1" {
		t.Errorf("Unexpected Request, %q", lines[13].Request)
	}

	if lines[26].Request[:7] != "OPTIONS" {
		t.Errorf("Unexpected Request, %q", lines[26].Request[:7])
	}

	if lines[32].HttpReferrer != "http://hvd-store.com/" {
		t.Errorf("Unexpected http referrer, %q", lines[32].HttpReferrer)
	}

	if lines[38].Request[:4] != "POST" {
		t.Errorf("Unexpected Request, %q", lines[38].Request[:4])
	}

	if lines[44].RemoteAddr != "141.212.122.170" {
		t.Errorf("Unexpected RemoteAddr, %q", lines[44].RemoteAddr)
	}
}
