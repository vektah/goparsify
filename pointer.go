package parsec

import (
	"strings"
	"unicode/utf8"
)

const (
	EOF rune = -1
)

func Input(s string) Pointer {
	return Pointer{s, 0}
}

type Pointer struct {
	input string
	pos   int
}

func (p Pointer) Advance(i int) Pointer {
	return Pointer{p.input, p.pos + i}
}

func (p Pointer) Get() string {
	return p.input[p.pos:]
}

func (p Pointer) Remaining() int {
	remaining := len(p.input) - p.pos
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (p Pointer) Next() (rune, Pointer) {
	if int(p.pos) >= len(p.input) {
		return EOF, p
	}
	r, w := utf8.DecodeRuneInString(p.input[p.pos:])
	return r, p.Advance(w)
}

func (p Pointer) HasPrefix(s string) bool {
	return strings.HasPrefix(p.input[p.pos:], s)
}

func (p Pointer) Accept(valid string) (string, Pointer) {
	c, newP := p.Next()
	if strings.ContainsRune(valid, c) {
		return string(c), newP
	}
	return "", p
}

func (p Pointer) AcceptRun(valid string) (string, Pointer) {
	matched := 0
	for p.pos+matched < len(p.input) {
		r, w := utf8.DecodeRuneInString(p.input[p.pos+matched:])
		if !strings.ContainsRune(valid, r) {
			break
		}
		matched += w
	}

	return p.input[p.pos : p.pos+matched], p.Advance(matched)
}

func (p Pointer) AcceptUntil(invalid string) (string, Pointer) {
	matched := 0
	for p.pos+matched < len(p.input) {
		r, w := utf8.DecodeRuneInString(p.input[p.pos+matched:])
		if strings.ContainsRune(invalid, r) {
			break
		}
		matched += w
	}

	return p.input[p.pos : p.pos+matched], p.Advance(matched)
}
