package goparsify

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Error struct {
	pos      int
	Expected string
}

func (e Error) Pos() int      { return e.pos }
func (e Error) Error() string { return fmt.Sprintf("offset %d: Expected %s", e.pos, e.Expected) }

type State struct {
	Input    string
	Pos      int
	Error    Error
	WSChars  string
	NoAutoWS bool
}

func (s *State) Advance(i int) {
	s.Pos += i
}

// AutoWS consumes all whitespace
func (s *State) AutoWS() {
	if s.NoAutoWS {
		return
	}
	s.WS()
}

func (s *State) WS() {
	for s.Pos < len(s.Input) {
		r, w := utf8.DecodeRuneInString(s.Input[s.Pos:])
		if !strings.ContainsRune(s.WSChars, r) {
			return
		}
		s.Pos += w

	}
}

func (s *State) Get() string {
	if s.Pos > len(s.Input) {
		return ""
	}
	return s.Input[s.Pos:]
}

func (s *State) ErrorHere(expected string) {
	s.Error.pos = s.Pos
	s.Error.Expected = expected
}

func (s *State) ClearError() {
	s.Error.Expected = ""
}

func (s *State) Errored() bool {
	return s.Error.Expected != ""
}

func InputString(input string) *State {
	return &State{Input: input, WSChars: "\t\n\v\f\r \x85\xA0"}
}
