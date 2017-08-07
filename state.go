package goparsify

import (
	"fmt"
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
	WSChars  []byte
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
loop:
	for s.Pos < len(s.Input) {
		// Pretty sure this is unicode safe as long as WSChars is only in the ascii range...
		for _, ws := range s.WSChars {
			if s.Input[s.Pos] == ws {
				s.Pos++
				continue loop
			}
		}

		return
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
	return &State{Input: input, WSChars: []byte("\t\n\v\f\r ")}
}
