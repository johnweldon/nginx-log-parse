package nginx

import (
	"fmt"
	"time"
)

// LogLine represents a generic line in the logfile, with no known structure.
type LogLine interface {
	fmt.Stringer
}

// RequestLine represents a standard 'combined' format log line in the nginx log.
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
