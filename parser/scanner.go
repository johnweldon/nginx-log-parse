package parser

import (
	"bufio"
	"bytes"
	"io"
)

type scanner struct {
	r         *bufio.Reader
	inQuote   bool
	inBracket bool
}

type token int

const (
	illegal token = iota
	eof
	eol
	space
	quote
	lbracket
	rbracket
	ident
)

var eofRune = rune(0)

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eofRune
	}
	return ch
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func (s *scanner) scan() (tok token, lit string) {
	ch := s.read()

	if isNewline(ch) {
		s.unread()
		return s.scanNewline()
	}
	switch ch {
	case eofRune:
		return eof, ""
	case ' ':
		return space, " "
	case '"':
		s.inQuote = !s.inQuote
		return quote, `"`
	case '[':
		s.inBracket = true
		return lbracket, "["
	case ']':
		s.inBracket = false
		return rbracket, "]"
	default:
		s.unread()
		return s.scanIdent()
	}
}

func (s *scanner) scanNewline() (tok token, lit string) {
	for {
		if ch := s.read(); ch == eofRune {
			break
		} else if !isNewline(ch) {
			s.unread()
			break
		}
	}
	return eol, "\n"
}

func (s *scanner) scanIdent() (tok token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eofRune {
			break
		} else if isNewline(ch) {
			s.unread()
			break
		} else if s.inQuote && ch != '"' {
			_, _ = buf.WriteRune(ch)
		} else if s.inBracket {
			if ch == ']' {
				s.unread()
				break
			}
			_, _ = buf.WriteRune(ch)
		} else if ch != ' ' && ch != '"' {
			_, _ = buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}
	return ident, buf.String()
}

func isNewline(ch rune) bool {
	switch ch {
	case '\r', '\n':
		return true
	default:
		return false
	}
}
