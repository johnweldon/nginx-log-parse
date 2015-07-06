package nginx

import "time"

type LogEntry struct {
	RemoteAddr    string
	RemoteUser    string
	TimeLocal     time.Time
	Request       string
	Status        int
	BodyBytesSent int
	HttpReferrer  string
	HttpUserAgent string
}
