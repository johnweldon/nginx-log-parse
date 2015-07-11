package util

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/johnweldon/nginx-log-parse/nginx"
)

type token struct {
	T Token
	L string
}

type handleToken struct {
	T Token
	F func(line nginx.LogLine, t token) error
}

type handleLine struct {
	N  string
	F  func() nginx.LogLine
	HT []handleToken
}

var supportedFormats = []handleLine{{
	N: "Combined", F: func() nginx.LogLine { return &nginx.LogEntry{} }, HT: []handleToken{
		// RemoteAddr
		{T: IDENT, F: setRemoteAddr},
		{T: SPACE, F: discard},
		// -
		{T: IDENT, F: discard},
		{T: SPACE, F: discard},
		// RemoteUser
		{T: IDENT, F: setRemoteUser},
		{T: SPACE, F: discard},
		// TimeLocal
		{T: LBRACKET, F: discard},
		{T: IDENT, F: setTimeLocal},
		{T: RBRACKET, F: discard},
		{T: SPACE, F: discard},
		// Request
		{T: QUOTE, F: discard},
		{T: IDENT, F: setRequest},
		{T: QUOTE, F: discard},
		{T: SPACE, F: discard},
		// Status
		{T: IDENT, F: setStatus},
		{T: SPACE, F: discard},
		// BodyBytesSent
		{T: IDENT, F: setBodyBytesSent},
		{T: SPACE, F: discard},
		// HttpReferrer
		{T: QUOTE, F: discard},
		{T: IDENT, F: setHttpReferrer},
		{T: QUOTE, F: discard},
		{T: SPACE, F: discard},
		// HttpUserAgent
		{T: QUOTE, F: discard},
		{T: IDENT, F: setHttpUserAgent},
		{T: QUOTE, F: discard},
	}}, {
	N: "Tail Divider", F: func() nginx.LogLine { return &nginx.DelimiterLine{} }, HT: []handleToken{
		// ==>
		{T: IDENT, F: discard},
		{T: SPACE, F: discard},
		// filename
		{T: IDENT, F: setDelimiterLine},
		//  <==
		{T: SPACE, F: discard},
		{T: IDENT, F: discard},
	}},
}

type Parser struct {
	s      *Scanner
	tokens []token
	LineCh chan nginx.LogLine
	tomb.Tomb
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{
		s:      NewScanner(r),
		LineCh: make(chan nginx.LogLine),
	}
	p.Go(p.loop)
	return p
}

func (p *Parser) GetRecords() []nginx.LogLine {
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

func (p *Parser) Stop() error {
	p.Kill(nil)
	return p.Wait()
}

func (p *Parser) loadLine() Token {
	p.tokens = []token{}
	for {
		if tok, lit := p.s.Scan(); tok != EOF && tok != EOL {
			p.tokens = append(p.tokens, token{T: tok, L: lit})
		} else {
			return tok
		}
	}
}

func (p *Parser) parseLine() nginx.LogLine {
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
	return &nginx.OtherEntry{Line: line}
}

func (p *Parser) loop() error {
	for {
		if tok := p.loadLine(); tok == EOF {
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

func discard(l nginx.LogLine, t token) error { return nil }

func setRemoteAddr(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.RemoteAddr = t.L
	return nil
}

func setRemoteUser(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.RemoteUser = t.L
	return nil
}

func setTimeLocal(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
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

func setRequest(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.Request = nginx.NewRequest(t.L)
	return nil
}

func setStatus(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
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

func setBodyBytesSent(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
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

func setHttpReferrer(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.HttpReferrer = t.L
	return nil
}

func setHttpUserAgent(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.LogEntry)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.HttpUserAgent = t.L
	return nil
}

func setDelimiterLine(l nginx.LogLine, t token) error {
	e, ok := l.(*nginx.DelimiterLine)
	if !ok {
		return fmt.Errorf("expected l to be a LogEntry")
	}
	e.Line = t.L
	return nil
}
