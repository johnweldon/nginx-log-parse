package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/johnweldon/nginx-log-parse/nginx"
)

type otherEntry struct {
	Line string
}

var _ nginx.LogLine = (*otherEntry)(nil)

func (o *otherEntry) String() string {
	return fmt.Sprintf(" unexpected line: %s", o.Line)
}

type delimiterLine struct {
	otherEntry
}

var _ nginx.LogLine = (*delimiterLine)(nil)

func (d *delimiterLine) String() string {
	return fmt.Sprintf("FILE: %s", d.Line)
}

type request struct {
	Method   string
	URI      string
	Protocol string
}

var _ nginx.RequestLine = (*logEntry)(nil)

func newRequest(ident string) request {
	raw := strings.Split(ident, " ")
	if len(raw) == 3 {
		return request{
			Method:   raw[0],
			URI:      raw[1],
			Protocol: raw[2],
		}
	}
	return request{}
}

type logEntry struct {
	RemoteAddr    string
	RemoteUser    string
	TimeLocal     time.Time
	Request       request
	Status        int
	BodyBytesSent int
	HTTPReferrer  string
	HTTPUserAgent string
}

func (e *logEntry) RequestHTTPReferrer() string  { return e.HTTPReferrer }
func (e *logEntry) RequestHTTPUserAgent() string { return e.HTTPUserAgent }
func (e *logEntry) RequestIP() string            { return e.RemoteAddr }
func (e *logEntry) RequestMethod() string        { return e.Request.Method }
func (e *logEntry) RequestProtocol() string      { return e.Request.Protocol }
func (e *logEntry) RequestTime() time.Time       { return e.TimeLocal }
func (e *logEntry) RequestURI() string           { return e.Request.URI }
func (e *logEntry) RequestUser() string          { return e.RemoteUser }
func (e *logEntry) ResponseBodyBytesSent() int   { return e.BodyBytesSent }
func (e *logEntry) ResponseStatus() int          { return e.Status }
func (e *logEntry) String() string {
	return fmt.Sprintf(
		"%20s %16s %d %4s %-30s %20s %s",
		e.TimeLocal.Format("2006-01-02 15:04:05"),
		e.RemoteAddr,
		e.Status,
		e.Request.Method,
		e.Request.URI,
		e.HTTPReferrer,
		e.HTTPUserAgent)
}
