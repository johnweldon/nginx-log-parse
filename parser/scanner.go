package parser

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	EOL
	SPACE
	QUOTE
	LBRACKET
	RBRACKET
	IDENT
)

var eof = rune(0)

type Scanner struct {
	r         *bufio.Reader
	inQuote   bool
	inBracket bool
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isNewline(ch) {
		s.unread()
		return s.scanNewline()
	}
	switch ch {
	case eof:
		return EOF, ""
	case ' ':
		return SPACE, " "
	case '"':
		s.inQuote = !s.inQuote
		return QUOTE, `"`
	case '[':
		s.inBracket = true
		return LBRACKET, "["
	case ']':
		s.inBracket = false
		return RBRACKET, "]"
	default:
		s.unread()
		return s.scanIdent()
	}
}

func (s *Scanner) scanNewline() (tok Token, lit string) {
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNewline(ch) {
			s.unread()
			break
		}
	}
	return EOL, "\n"
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
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
	return IDENT, buf.String()
}

func isNewline(ch rune) bool {
	switch ch {
	case '\r', '\n':
		return true
	default:
		return false
	}
}
