package parser

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/johnweldon/nginx-log-parse/nginx"
)

// Engine describes the minimal parser engine for the nginx logfiles.
type Engine interface {
	GetRecords() []nginx.LogLine
	LogLines() <-chan nginx.LogLine
	Dying() <-chan struct{}
	Stop() error
}

type logFileParser struct {
	s      *scanner
	tokens []tokenLiteral
	LineCh chan nginx.LogLine
	tomb.Tomb
}

type tokenLiteral struct {
	T token
	L string
}

type handleToken struct {
	T token
	F func(line nginx.LogLine, t tokenLiteral) error
}

type handleLine struct {
	N  string
	F  func() nginx.LogLine
	HT []handleToken
}

var supportedFormats = []handleLine{{
	N: "Combined", F: func() nginx.LogLine { return &logEntry{} }, HT: []handleToken{
		// RemoteAddr
		{T: ident, F: setRemoteAddr},
		{T: space, F: discard},
		// -
		{T: ident, F: discard},
		{T: space, F: discard},
		// RemoteUser
		{T: ident, F: setRemoteUser},
		{T: space, F: discard},
		// TimeLocal
		{T: lbracket, F: discard},
		{T: ident, F: setTimeLocal},
		{T: rbracket, F: discard},
		{T: space, F: discard},
		// Request
		{T: quote, F: discard},
		{T: ident, F: setRequest},
		{T: quote, F: discard},
		{T: space, F: discard},
		// Status
		{T: ident, F: setStatus},
		{T: space, F: discard},
		// BodyBytesSent
		{T: ident, F: setBodyBytesSent},
		{T: space, F: discard},
		// HTTPReferrer
		{T: quote, F: discard},
		{T: ident, F: setHTTPReferrer},
		{T: quote, F: discard},
		{T: space, F: discard},
		// HTTPUserAgent
		{T: quote, F: discard},
		{T: ident, F: setHTTPUserAgent},
		{T: quote, F: discard},
	}}, {
	N: "Tail Divider", F: func() nginx.LogLine { return &delimiterLine{} }, HT: []handleToken{
		// ==> (or any contiguous symbol actually)
		{T: ident, F: discard},
		{T: space, F: discard},
		// filename
		{T: ident, F: setDelimiterLine},
		//  <== (or any contiguous symbol actually)
		{T: space, F: discard},
		{T: ident, F: discard},
	}},
}

// NewEngine returns an Engine implementation.
func NewEngine(r io.Reader) Engine { return newLogFileParser(r) }

func (p *logFileParser) GetRecords() []nginx.LogLine {
	var records []nginx.LogLine
	for {
		select {
		case line := <-p.LineCh:
			record, ok := line.(nginx.LogLine)
			if !ok {
				continue
			}
			records = append(records, record)
		case <-p.Dying():
			return records
		}
	}
}

func (p *logFileParser) LogLines() <-chan nginx.LogLine {
	return p.LineCh
}

func (p *logFileParser) Stop() error {
	p.Kill(nil)
	return p.Wait()
}

func newLogFileParser(r io.Reader) *logFileParser {
	p := &logFileParser{
		s:      newScanner(r),
		LineCh: make(chan nginx.LogLine),
	}
	p.Go(p.loop)
	return p
}

func (p *logFileParser) loadLine() token {
	p.tokens = []tokenLiteral{}
	for {
		if tok, lit := p.s.scan(); tok != eof && tok != eol {
			p.tokens = append(p.tokens, tokenLiteral{T: tok, L: lit})
		} else {
			return tok
		}
	}
}

func (p *logFileParser) parseLine() nginx.LogLine {
	for _, handler := range supportedFormats {
		if len(p.tokens) != len(handler.HT) {
			continue
		}
		l := handler.F()
		for ix, tok := range p.tokens {
			if tok.T != handler.HT[ix].T {
				continue
			}
			if err := handler.HT[ix].F(l, tok); err != nil {
				continue
			}
		}
		return l
	}
	line := ""
	for _, tok := range p.tokens {
		line = line + tok.L
	}
	return &otherEntry{Line: line}
}

func (p *logFileParser) loop() error {
	for {
		if tok := p.loadLine(); tok == eof {
			close(p.LineCh)
			return nil
		}

		line := p.parseLine()

		select {
		case p.LineCh <- line:
		case <-p.Dying():
			close(p.LineCh)
			return nil
		}
	}
}

func discard(l nginx.LogLine, t tokenLiteral) error { return nil }

func setRemoteAddr(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.RemoteAddr = t.L
	return nil
}

func setRemoteUser(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.RemoteUser = t.L
	return nil
}

func setTimeLocal(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	ti, err := time.Parse("2/Jan/2006:15:04:05 -0700", t.L)
	if err != nil {
		return fmt.Errorf("expected time local, got %q: %v", t.L, err)
	}
	e.TimeLocal = ti
	return nil
}

func setRequest(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.Request = newRequest(t.L)
	return nil
}

func setStatus(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	status, err := strconv.ParseInt(t.L, 10, 0)
	if err != nil {
		return fmt.Errorf("expected Status, got %q: %v", t.L, err)
	}
	e.Status = int(status)
	return nil
}

func setBodyBytesSent(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	status, err := strconv.ParseInt(t.L, 10, 0)
	if err != nil {
		return fmt.Errorf("expected BodyBytesSent, got %q: %v", t.L, err)
	}
	e.BodyBytesSent = int(status)
	return nil
}

func setHTTPReferrer(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.HTTPReferrer = t.L
	return nil
}

func setHTTPUserAgent(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*logEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.HTTPUserAgent = t.L
	return nil
}

func setDelimiterLine(l nginx.LogLine, t tokenLiteral) error {
	e, ok := l.(*delimiterLine)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.Line = t.L
	return nil
}
