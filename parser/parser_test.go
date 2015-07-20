package parser_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/johnweldon/nginx-log-parse/nginx"
	. "github.com/johnweldon/nginx-log-parse/parser"
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

func isLogLine(l nginx.LogLine) bool { return true }
func isRequestLine(l nginx.LogLine) bool {
	_, ok := l.(nginx.RequestLine)
	return ok
}

var (
	fnLine           = func(l nginx.LogLine) interface{} { return l.String() }
	fnRequestIP      = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestIP() }
	fnResponseStatus = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).ResponseStatus() }
	fnResponseBytes  = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).ResponseBodyBytesSent() }
	fnUserAgent      = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestHTTPUserAgent() }
	fnURI            = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestURI() }
	fnRequestMethod  = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestMethod() }
	fnReferrer       = func(l nginx.LogLine) interface{} { return l.(nginx.RequestLine).RequestHTTPReferrer() }
)

func TestParse(t *testing.T) {
	tests := []testCases{
		{
			Lines:  46,
			Source: goodLog,
			Tests: []testCase{
				{
					Index:        0,
					ExpectedType: isLogLine,
					Got:          fnLine,
					Expect:       " unexpected line: ",
				},
				{
					Index:        1,
					ExpectedType: isRequestLine,
					Got:          fnRequestIP,
					Expect:       "10.10.10.15",
				},
				{
					Index:        5,
					ExpectedType: isRequestLine,
					Got:          fnResponseStatus,
					Expect:       404,
				},
				{
					Index:        8,
					ExpectedType: isRequestLine,
					Got:          func(l nginx.LogLine) interface{} { return fnUserAgent(l).(string)[13:22] },
					Expect:       "Macintosh",
				},
				{
					Index:        14,
					ExpectedType: isRequestLine,
					Got:          fnURI,
					Expect:       "/filter/tips",
				},
				{
					Index:        27,
					ExpectedType: isRequestLine,
					Got:          fnRequestMethod,
					Expect:       "OPTIONS",
				},
				{
					Index:        33,
					ExpectedType: isRequestLine,
					Got:          fnReferrer,
					Expect:       "http://hvd-store.com/",
				},
				{
					Index:        39,
					ExpectedType: isRequestLine,
					Got:          fnRequestMethod,
					Expect:       "POST",
				},
				{
					Index:        45,
					ExpectedType: isRequestLine,
					Got:          fnRequestIP,
					Expect:       "141.212.122.170",
				},
			},
		}, {
			Lines:  11,
			Source: mixedLog,
			Tests: []testCase{
				{
					Index:        0,
					ExpectedType: isLogLine,
					Got:          fnLine,
					Expect:       " unexpected line: ",
				},
				{
					Index:        1,
					ExpectedType: isRequestLine,
					Got:          fnURI,
					Expect:       "/",
				},
				{
					Index:        5,
					ExpectedType: isLogLine,
					Got:          fnLine,
					Expect:       " unexpected line: asdf",
				},
				{
					Index:        6,
					ExpectedType: isLogLine,
					Got:          fnLine,
					Expect:       "FILE: another.file",
				},
				{
					Index:        7,
					ExpectedType: isRequestLine,
					Got:          fnResponseBytes,
					Expect:       181,
				},
				{
					Index:        10,
					ExpectedType: isRequestLine,
					Got:          fnRequestIP,
					Expect:       "113.247.42.246",
				},
			},
		},
	}

	for tci, testCase := range tests {
		p := NewLogFileParser(bytes.NewReader([]byte(testCase.Source)))
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
				t.Errorf("testCase: %d: unexpected type for test %d: %T", tci, ti, line)
				continue
			}
			if !reflect.DeepEqual(test.Got(line), test.Expect) {
				t.Errorf("testCase: %d: test %d: expected %v, got %v", tci, ti, test.Expect, test.Got(line))
				continue
			}

		}
	}
}
