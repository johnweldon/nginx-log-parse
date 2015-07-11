package util_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/johnweldon/nginx-log-parse/nginx"
	. "github.com/johnweldon/nginx-log-parse/util"
)

type testCase struct {
	Index        int
	ExpectedType func(nginx.LogLine) bool
	Got          func(nginx.LogLine) interface{}
	Expect       interface{}
}
type testCases struct {
	Lines  int
	Source string
	Tests  []testCase
}

func isRequestLine(l nginx.LogLine) bool {
	_, ok := l.(nginx.RequestLine)
	return ok
}

var (
	fnRequestIP      = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestIP() }
	fnResponseStatus = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).ResponseStatus() }
	fnResponseBytes  = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).ResponseBodyBytesSent() }
	fnUserAgent      = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestHttpUserAgent() }
	fnURI            = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestURI() }
	fnRequestMethod  = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestMethod() }
	fnReferrer       = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestHttpReferrer() }
)

func TestParse(t *testing.T) {
	tests := []testCases{
		{
			Lines:  45,
			Source: goodLog,
			Tests: []testCase{
				{
					Index:        0,
					ExpectedType: isRequestLine,
					Got:          fnRequestIP,
					Expect:       "10.10.10.15",
				},
				{
					Index:        4,
					ExpectedType: isRequestLine,
					Got:          fnResponseStatus,
					Expect:       404,
				},
				{
					Index:        7,
					ExpectedType: isRequestLine,
					Got:          func(l nginx.LogLine) interface{} { return fnUserAgent(l).(string)[13:22] },
					Expect:       "Macintosh",
				},
				{
					Index:        13,
					ExpectedType: isRequestLine,
					Got:          fnURI,
					Expect:       "/filter/tips",
				},
				{
					Index:        26,
					ExpectedType: isRequestLine,
					Got:          fnRequestMethod,
					Expect:       "OPTIONS",
				},
				{
					Index:        32,
					ExpectedType: isRequestLine,
					Got:          fnReferrer,
					Expect:       "http://hvd-store.com/",
				},
				{
					Index:        38,
					ExpectedType: isRequestLine,
					Got:          fnRequestMethod,
					Expect:       "POST",
				},
				{
					Index:        44,
					ExpectedType: isRequestLine,
					Got:          fnRequestIP,
					Expect:       "141.212.122.170",
				},
			},
		}, {
			Lines:  8,
			Source: mixedLog,
			Tests: []testCase{
				{
					Index:        0,
					ExpectedType: isRequestLine,
					Got:          fnURI,
					Expect:       "/",
				},
				{
					Index:        7,
					ExpectedType: isRequestLine,
					Got:          fnResponseBytes,
					Expect:       396,
				},
			},
		},
	}

	for tci, testCase := range tests {
		p := NewParser(bytes.NewReader([]byte(testCase.Source)))
		lines := p.GetRecords()
		if len(lines) != testCase.Lines {
			t.Fatalf("expected %d lines, but got %d", testCase.Lines, len(lines))
		}
		for ti, test := range testCase.Tests {
			line := lines[test.Index]
			if line == nil {
				t.Errorf("testCase: %d: test %d: line is nil", tci, ti)
				continue
			}
			if !test.ExpectedType(line) {
				t.Errorf("testCase: %d: unexpected type for test %d: %t", tci, ti, line)
				continue
			}
			if !reflect.DeepEqual(test.Got(line), test.Expect) {
				t.Errorf("testCase: %d: test %d: expected %v, got %v", tci, ti, test.Expect, test.Got(line))
				continue
			}

		}
	}
}
