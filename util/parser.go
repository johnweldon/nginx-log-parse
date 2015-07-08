package util

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"gopkg.in/tomb.v2"

	"github.com/johnweldon/nginx-log-parse/nginx"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
	EntryCh chan *nginx.LogEntry
	Log     chan string
	tomb.Tomb
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{
		s:       NewScanner(r),
		EntryCh: make(chan *nginx.LogEntry),
		Log:     make(chan string),
	}
	p.Go(p.loop)
	return p
}

func (p *Parser) loop() error {
	for {
		line, err := p.Parse()
		if err != nil {
			if p.IsEOF() {
				p.closeChannels()
				return err
			}
			p.Log <- fmt.Sprintf("Warning: %q", err)
		}
		if line == nil {
			continue
		}
		select {
		case p.EntryCh <- line:
		case <-p.Dying():
			p.closeChannels()
			return nil
		}
	}
}

func (p *Parser) Stop() error {
	p.Kill(nil)
	return p.Wait()
}

func (p *Parser) closeChannels() {
	close(p.EntryCh)
	close(p.Log)
}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	tok, lit = p.s.Scan()
	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *Parser) discardLine() {
	if p.buf.tok == EOF || p.buf.tok == EOL {
		return
	}
	p.buf.n = 0
	for {
		tok, _ := p.s.Scan()
		if tok == EOF || tok == EOL {
			return
		}
	}
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) IsEOF() bool {
	return p.buf.tok == EOF
}

type assign func(string) error

func (p *Parser) ident(expect string, fn assign) error {
	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return fmt.Errorf("expected %s, got %q", expect, lit)
	} else {
		if fn != nil {
			if err := fn(lit); err != nil {
				return err
			}
		}
		return nil
	}
}

func (p *Parser) space() error {
	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return fmt.Errorf("expected space, got %q", lit)
	}
	return nil
}

func (p *Parser) lbracket() error {
	if tok, lit := p.scan(); tok != LBRACKET {
		p.discardLine()
		return fmt.Errorf("expected [, got %q", lit)
	}
	return nil
}

func (p *Parser) rbracket() error {
	if tok, lit := p.scan(); tok != RBRACKET {
		p.discardLine()
		return fmt.Errorf("expected ], got %q", lit)
	}
	return nil
}

func (p *Parser) quote() error {
	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return fmt.Errorf("expected quote, got %q", lit)
	}
	return nil
}

func (p *Parser) eol() error {
	if tok, lit := p.scan(); tok != EOL {
		p.discardLine()
		return fmt.Errorf("expected eol, got %q", lit)
	}
	return nil
}

func (p *Parser) Parse() (*nginx.LogEntry, error) {
	line := &nginx.LogEntry{}

	if err := p.ident("RemoteAddr", func(i string) error { line.RemoteAddr = i; return nil }); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.ident("-", nil); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.ident("RemoteUser", func(i string) error { line.RemoteUser = i; return nil }); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.lbracket(); err != nil {
		return nil, err
	}

	if err := p.ident("TimeLocal", func(i string) error {
		t, err := time.Parse("2/Jan/2006:15:04:05 -0700", i)
		if err != nil {
			p.discardLine()
			return fmt.Errorf("expected time local, got %q: %v", i, err)
		}
		line.TimeLocal = t
		return nil
	}); err != nil {
		return nil, err
	}

	if err := p.rbracket(); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.ident("Request", func(i string) error { line.Request = nginx.NewRequest(i); return nil }); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.ident("Status", func(i string) error {
		status, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			p.discardLine()
			return fmt.Errorf("expected Status, got %q: %v", i, err)
		}
		line.Status = int(status)
		return nil
	}); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.ident("BodyBytesSent", func(i string) error {
		sent, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			p.discardLine()
			return fmt.Errorf("expected BodyBytesSent, got %q: %v", i, err)
		}
		line.BodyBytesSent = int(sent)
		return nil
	}); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.ident("HttpReferrer", func(i string) error { line.HttpReferrer = i; return nil }); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.space(); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.ident("HttpUserAgent", func(i string) error { line.HttpUserAgent = i; return nil }); err != nil {
		return nil, err
	}

	if err := p.quote(); err != nil {
		return nil, err
	}

	if err := p.eol(); err != nil {
		return nil, err
	}

	return line, nil
}

func (p *Parser) GetRecords() []*nginx.LogEntry {
	var records []*nginx.LogEntry
	for {
		select {
		case record := <-p.EntryCh:
			if record == nil {
				continue
			}
			records = append(records, record)
		case <-p.Log:
			continue
		case <-p.Dying():
			return records
		}
	}
}
