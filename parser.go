package main

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
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

func (p *Parser) Parse() (*LogLine, error) {
	line := &LogLine{}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected ip address, got %q", lit)
	} else {
		line.RemoteAddr = lit
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT || lit != "-" {
		p.discardLine()
		return nil, fmt.Errorf("expected -, got %q", lit)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected remote user, got %q", lit)
	} else {
		line.RemoteUser = lit
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != LBRACKET {
		p.discardLine()
		return nil, fmt.Errorf("expected [, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected time local, got %q", lit)
	} else {
		t, err := time.Parse("2/Jan/2006:15:04:05 -0700", lit)
		if err != nil {
			p.discardLine()
			return nil, fmt.Errorf("expected time local, got %q: %v", lit, err)
		}
		line.TimeLocal = t
	}

	if tok, lit := p.scan(); tok != RBRACKET {
		p.discardLine()
		return nil, fmt.Errorf("expected ], got %q", lit)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected request, got %q", lit)
	} else {
		line.Request = lit
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected status, got %q", lit)
	} else {
		status, err := strconv.ParseInt(lit, 10, 0)
		if err != nil {
			p.discardLine()
			return nil, fmt.Errorf("expected status, got %q: %v", lit, err)
		}
		line.Status = int(status)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected body bytes sent, got %q", lit)
	} else {
		sent, err := strconv.ParseInt(lit, 10, 0)
		if err != nil {
			p.discardLine()
			return nil, fmt.Errorf("expected body bytes sent, got %q: %v", lit, err)
		}
		line.BodyBytesSent = int(sent)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected http referrer, got %q", lit)
	} else {
		line.HttpReferrer = lit
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != SPACE {
		p.discardLine()
		return nil, fmt.Errorf("expected space, got %q", lit)
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != IDENT {
		p.discardLine()
		return nil, fmt.Errorf("expected http user agent, got %q", lit)
	} else {
		line.HttpUserAgent = lit
	}

	if tok, lit := p.scan(); tok != QUOTE {
		p.discardLine()
		return nil, fmt.Errorf("expected quote, got %q", lit)
	}

	if tok, lit := p.scan(); tok != EOL {
		p.discardLine()
		return nil, fmt.Errorf("expected eol, got %q", lit)
	}

	return line, nil
}

func (p *Parser) GetRecords() []LogLine {
	lines := []LogLine{}
	for {
		if line, err := p.Parse(); err != nil {
			if p.IsEOF() {
				break
			}
			continue
		} else {
			lines = append(lines, *line)
		}
	}
	return lines
}
