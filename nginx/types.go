package nginx

import (
	"fmt"
	"strings"
	"time"
)

type LogLine interface {
	fmt.Stringer
}

var _ LogLine = (*OtherEntry)(nil)
var _ LogLine = (*DelimiterLine)(nil)

type RequestLine interface {
	LogLine
	RequestHTTPReferrer() string
	RequestHTTPUserAgent() string
	RequestIP() string
	RequestMethod() string
	RequestProtocol() string
	RequestTime() time.Time
	RequestURI() string
	RequestUser() string
	ResponseBodyBytesSent() int
	ResponseStatus() int
}

var _ RequestLine = (*LogEntry)(nil)

type OtherEntry struct {
	Line string
}

func (o *OtherEntry) String() string {
	return fmt.Sprintf(" unexpected line: %s", o.Line)
}

type DelimiterLine struct {
	OtherEntry
}

func (d *DelimiterLine) String() string {
	return fmt.Sprintf("FILE: %s", d.Line)
}

type Request struct {
	Method   string
	URI      string
	Protocol string
}

func NewRequest(ident string) Request {
	parts := strings.Split(ident, " ")
	if len(parts) == 3 {
		return Request{
			Method:   parts[0],
			URI:      parts[1],
			Protocol: parts[2],
		}
	}
	return Request{}
}

type LogEntry struct {
	RemoteAddr    string
	RemoteUser    string
	TimeLocal     time.Time
	Request       Request
	Status        int
	BodyBytesSent int
	HTTPReferrer  string
	HTTPUserAgent string
}

func (e *LogEntry) RequestHTTPReferrer() string  { return e.HTTPReferrer }
func (e *LogEntry) RequestHTTPUserAgent() string { return e.HTTPUserAgent }
func (e *LogEntry) RequestIP() string            { return e.RemoteAddr }
func (e *LogEntry) RequestMethod() string        { return e.Request.Method }
func (e *LogEntry) RequestProtocol() string      { return e.Request.Protocol }
func (e *LogEntry) RequestTime() time.Time       { return e.TimeLocal }
func (e *LogEntry) RequestURI() string           { return e.Request.URI }
func (e *LogEntry) RequestUser() string          { return e.RemoteUser }
func (e *LogEntry) ResponseBodyBytesSent() int   { return e.BodyBytesSent }
func (e *LogEntry) ResponseStatus() int          { return e.Status }
func (e *LogEntry) String() string {
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
